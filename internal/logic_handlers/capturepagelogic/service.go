package capturepagelogic

import (
	"context"
	"fmt"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/sysconsts"
	"github.com/dembygenesis/local.tools/internal/utilities/errs"
	"github.com/dembygenesis/local.tools/internal/utilities/validationutils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Config struct {
	TxProvider persistence.TransactionProvider `json:"tx_provider" validate:"required"`
	Logger     *logrus.Entry                   `json:"Logger" validate:"required"`
	Persistor  persistor                       `json:"Persistor" validate:"required"`
}

func (i *Config) Validate() error {
	return validationutils.Validate(i)
}

type Service struct {
	cfg *Config
}

func New(cfg *Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &Service{cfg}, nil
}

func (i *Service) CreateCapturePage(ctx context.Context, params *model.CreateCapturePage) (*model.CapturePage, error) {
	if err := params.Validate(); err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("validate: %w", err),
		})
	}

	tx, err := i.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %w", err),
		})
	}
	defer tx.Rollback(ctx)

	capturePage := params.ToCapturePage()

	exists, err := i.cfg.Persistor.GetCapturePageByName(ctx, tx, capturePage.Name)
	if err != nil {
		if !strings.Contains(err.Error(), sysconsts.ErrExpectedExactlyOneEntry) {
			return nil, errs.New(&errs.Cfg{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("check capture page unique: %w", err),
			})
		}
	}
	if exists != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("capture page already exists"),
		})
	}

	capturePageCreated, err := i.cfg.Persistor.CreateCapturePage(ctx, tx, capturePage)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("create: %w", err),
		})
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("commit: %v", err),
		})
	}

	return capturePageCreated, nil
}

func (i *Service) AddCapturePage(ctx context.Context, params *model.CreateCapturePage) (*model.CapturePage, error) {
	if err := params.Validate(); err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("validate: %w", err),
		})
	}

	tx, err := i.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %w", err),
		})
	}
	defer tx.Rollback(ctx)

	exists, err := i.cfg.Persistor.GetCapturePageByName(ctx, tx, params.Name)
	if err != nil {
		if !strings.Contains(err.Error(), sysconsts.ErrExpectedExactlyOneEntry) {
			return nil, errs.New(&errs.Cfg{
				StatusCode: http.StatusBadRequest,
				Err:        fmt.Errorf("check capture page unique: %v", err),
			})
		}
	}
	if exists != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("capture page already exists"),
		})
	}

	capturePage, err := i.cfg.Persistor.AddCapturePage(ctx, tx, params)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("create: %w", err),
		})
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("commit: %w", err),
		})
	}

	return capturePage, nil
}

func (i *Service) ListCapturePages(
	ctx context.Context,
	filter *model.CapturePageFilters,
) (*model.PaginatedCapturePages, error) {
	db, err := i.cfg.TxProvider.Db(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %v", err),
		})
	}

	paginated, err := i.cfg.Persistor.GetCapturePages(ctx, db, filter)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get capture pages: %v", err),
		})
	}

	return paginated, nil
}

func (i *Service) UpdateCapturePage(ctx context.Context, params *model.UpdateCapturePage) (*model.CapturePage, error) {
	tx, err := i.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return nil, errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %w", err),
		})
	}
	defer tx.Rollback(ctx)

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	capturePage, err := i.cfg.Persistor.UpdateCapturePage(ctx, tx, params)
	if err != nil {
		return nil, fmt.Errorf("update capture page: %w", err)
	}

	return capturePage, nil
}

func (s *Service) DeleteCapturePage(ctx context.Context, params *model.DeleteCapturePage) error {
	tx, err := s.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %w", err),
		})
	}
	defer tx.Rollback(ctx)

	err = s.cfg.Persistor.DeleteCapturePage(ctx, tx, params.ID)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("delete capture page: %w", err),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("commit transaction: %w", err),
		})
	}

	return nil
}

func (s *Service) RestoreCapturePage(ctx context.Context, params *model.RestoreCapturePage) error {
	tx, err := s.cfg.TxProvider.Tx(ctx)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("get db: %w", err),
		})
	}
	defer tx.Rollback(ctx)

	err = s.cfg.Persistor.RestoreCapturePage(ctx, tx, params.ID)
	if err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("restore capture page: %w", err),
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return errs.New(&errs.Cfg{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("commit transaction: %w", err),
		})
	}

	return nil
}
