package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (m *Repository) GetCapturePageSetById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.CapturePageSet, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getCapturePageSet(ctx, ctxExec, &model.CapturePageSetFilters{IdsIn: []int{id}})
	if err != nil {
		return nil, fmt.Errorf("read capture pages: %v", err)
	}

	if res.Pagination.RowCount != 1 {
		return nil, fmt.Errorf(sysconsts.ErrExpectedExactlyOneEntry, mysqlmodel.TableNames.CapturePageSet)
	}

	return &res.CapturePages[0], nil
}

func (m *Repository) GetCapturePageSet(ctx context.Context, tx persistence.TransactionHandler, filters *model.CapturePageSetFilters) (*model.PaginatedCapturePageSet, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %v", err)
	}

	res, err := m.getCapturePageSet(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read capture pages: %v", err)
	}

	return res, nil
}

func (m *Repository) getCapturePageSet(
	ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.CapturePageSetFilters,
) (*model.PaginatedCapturePageSet, error) {
	var (
		paginated  model.PaginatedCapturePageSet
		pagination = model.NewPagination()
		res        = make([]model.CapturePageSet, 0)
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := make([]qm.QueryMod, 0)

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.CapturePageSetWhere.ID.IN(filters.IdsIn))
		}
	}

	q := mysqlmodel.CapturePageSets(queryMods...)
	qCount, err := q.Count(ctx, ctxExec)
	if err != nil {
		return nil, fmt.Errorf("get capture pages count: %v", err)
	}

	page := pagination.Page
	maxRows := pagination.MaxRows
	if filters != nil {
		if filters.Page.Valid {
			page = filters.Page.Int
		}
		if filters.MaxRows.Valid {
			maxRows = filters.MaxRows.Int
		}
	}

	pagination.SetQueryBoundaries(page, maxRows, int(qCount))

	queryMods = append(queryMods, qm.Limit(pagination.MaxRows), qm.Offset(pagination.Offset))
	q = mysqlmodel.CapturePageSets(queryMods...)

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get capture pages: %v", err)
	}

	pagination.RowCount = len(res)
	paginated.CapturePages = res
	paginated.Pagination = pagination

	return &paginated, nil
}
