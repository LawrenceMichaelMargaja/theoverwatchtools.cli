package mysqlstore

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"strconv"
)

func (m *Repository) GetClickTrackerById(ctx context.Context, tx persistence.TransactionHandler, id int) (*model.ClickTracker, error) {
	paginated, err := m.GetClickTrackers(ctx, tx, &model.ClickTrackerFilters{
		IdsIn: []int{id},
	})

	if err != nil {
		return nil, fmt.Errorf("click tracker filtered by id: %w", err)
	}

	if paginated.Pagination.RowCount != 1 {
		return nil, fmt.Errorf(sysconsts.ErrExpectedExactlyOneEntry, id)
	}

	return &paginated.ClickTrackers[0], nil
}

func (m *Repository) GetClickTrackerByName(ctx context.Context, tx persistence.TransactionHandler, name string) (*model.ClickTracker, error) {
	paginated, err := m.GetClickTrackers(ctx, tx, &model.ClickTrackerFilters{
		ClickTrackerNameIn: []string{name},
	})
	if err != nil {
		return nil, fmt.Errorf("click tracker filtered by name: %w", err)
	}

	if paginated.Pagination.RowCount != 1 {
		return nil, errors.New(sysconsts.ErrExpectedExactlyOneEntry)
	}

	return &paginated.ClickTrackers[0], nil
}

func (m *Repository) GetClickTrackers(ctx context.Context, tx persistence.TransactionHandler, filters *model.ClickTrackerFilters) (*model.PaginatedClickTrackers, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	res, err := m.getClickTrackers(ctx, ctxExec, filters)
	if err != nil {
		return nil, fmt.Errorf("read click tracker: %w", err)
	}

	return res, nil
}

func (m *Repository) getClickTrackers(ctx context.Context,
	ctxExec boil.ContextExecutor,
	filters *model.ClickTrackerFilters) (*model.PaginatedClickTrackers, error) {

	var (
		paginated  model.PaginatedClickTrackers
		pagination = model.NewPagination()
		res        = make([]model.ClickTracker, 0)
		err        error
	)

	ctx, cancel := context.WithTimeout(ctx, m.cfg.QueryTimeouts.Query)
	defer cancel()

	queryMods := []qm.QueryMod{
		qm.InnerJoin(
			fmt.Sprintf(
				"%s ON %s.%s = %s.%s",
				mysqlmodel.TableNames.User,
				mysqlmodel.TableNames.User,
				mysqlmodel.UserColumns.ID,
				mysqlmodel.TableNames.ClickTracker,
				mysqlmodel.ClickTrackerColumns.CreatedBy,
			),
		),
		qm.Select(
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTracker,
				mysqlmodel.ClickTrackerColumns.ID,
				mysqlmodel.ClickTrackerColumns.ID,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTracker,
				mysqlmodel.ClickTrackerColumns.Name,
				mysqlmodel.ClickTrackerColumns.Name,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTracker,
				mysqlmodel.ClickTrackerColumns.CreatedBy,
				mysqlmodel.ClickTrackerColumns.CreatedBy,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTracker,
				mysqlmodel.ClickTrackerColumns.LastUpdatedBy,
				mysqlmodel.ClickTrackerColumns.LastUpdatedBy,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTracker,
				mysqlmodel.ClickTrackerColumns.CreatedAt,
				mysqlmodel.ClickTrackerColumns.CreatedAt,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTracker,
				mysqlmodel.ClickTrackerColumns.LastUpdatedAt,
				mysqlmodel.ClickTrackerColumns.LastUpdatedAt,
			),
			fmt.Sprintf("%s.%s AS %s",
				mysqlmodel.TableNames.ClickTracker,
				mysqlmodel.ClickTrackerColumns.IsActive,
				mysqlmodel.ClickTrackerColumns.IsActive,
			),
		),
	}

	if filters != nil {
		if len(filters.IdsIn) > 0 {
			fmt.Println("the filters --- ", filters.IdsIn)
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.ID.IN(filters.IdsIn))
		}

		if filters.ClickTrackerIsActive.Valid {
			fmt.Println("the filters valid --- ", filters.ClickTrackerIsActive)
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.IsActive.EQ(filters.ClickTrackerIsActive.Bool))
		}

		if len(filters.ClickTrackerNameIn) > 0 {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.Name.IN(filters.ClickTrackerNameIn))
		}

		if filters.CreatedBy.Valid {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.CreatedBy.EQ(filters.CreatedBy))
		}

		if filters.LastUpdatedBy.Valid {
			queryMods = append(queryMods, mysqlmodel.ClickTrackerWhere.LastUpdatedBy.EQ(filters.LastUpdatedBy))
		}
	}

	q := mysqlmodel.ClickTrackers(queryMods...)
	totalCount, err := q.Count(ctx, ctxExec)
	if err != nil {
		return nil, fmt.Errorf("get click tracker count: %w", err)
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

	pagination.SetQueryBoundaries(page, maxRows, int(totalCount))

	queryMods = append(queryMods, qm.Limit(pagination.MaxRows), qm.Offset(pagination.Offset))
	q = mysqlmodel.ClickTrackers(queryMods...)

	if err = q.Bind(ctx, ctxExec, &res); err != nil {
		return nil, fmt.Errorf("get click tracker: %w", err)
	}

	pagination.RowCount = len(res)
	paginated.ClickTrackers = res
	paginated.Pagination = pagination

	return &paginated, nil
}

func (m *Repository) CreateClickTracker(ctx context.Context, tx persistence.TransactionHandler, clickTracker *model.ClickTracker) (*model.ClickTracker, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	createdBY, _ := strconv.Atoi(clickTracker.CreatedBy)
	lastUpdateBY, _ := strconv.Atoi(clickTracker.LastUpdatedBy)
	entry := mysqlmodel.ClickTracker{
		Name:          clickTracker.Name,
		CreatedBy:     null.IntFrom(createdBY),
		LastUpdatedBy: null.IntFrom(lastUpdateBY),
	}
	if err = entry.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert click tracker: %w", err)
	}

	clickTracker, err = m.GetClickTrackerById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get click tracker by id: %w", err)
	}

	return clickTracker, nil
}

func (m *Repository) AddClickTracker(ctx context.Context, tx persistence.TransactionHandler, clickTracker *model.CreateClickTracker) (*model.ClickTracker, error) {
	ctxExec, err := mysqltx.GetCtxExecutor(tx)
	if err != nil {
		return nil, fmt.Errorf("extract context executor: %w", err)
	}

	entry := &mysqlmodel.ClickTracker{
		Name:          clickTracker.Name,
		CreatedBy:     null.IntFrom(clickTracker.UserId),
		LastUpdatedBy: null.IntFrom(clickTracker.UserId),
	}
	if err = entry.Insert(ctx, ctxExec, boil.Infer()); err != nil {
		return nil, fmt.Errorf("insert click tracker: %w", err)
	}

	createClickTracker, err := m.GetClickTrackerById(ctx, tx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("get click tracker by id: %w", err)
	}

	return createClickTracker, nil
}
