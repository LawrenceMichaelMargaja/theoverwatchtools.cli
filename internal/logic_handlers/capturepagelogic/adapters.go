package capturepagelogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . persistor
type persistor interface {
	GetCapturePages(ctx context.Context, tx persistence.TransactionHandler, filters *model.CapturePageFilters) (*model.PaginatedCapturePages, error)
	CreateCapturePage(ctx context.Context, tx persistence.TransactionHandler, capturePage *model.CapturePage) (*model.CapturePage, error)
	AddCapturePage(ctx context.Context, tx persistence.TransactionHandler, capturePage *model.CreateCapturePage) (*model.CapturePage, error)
	GetCapturePageByName(ctx context.Context, tx persistence.TransactionHandler, name string) (*model.CapturePage, error)
	UpdateCapturePage(ctx context.Context, tx persistence.TransactionHandler, params *model.UpdateCapturePage) (*model.CapturePage, error)
	DeleteCapturePage(ctx context.Context, tx persistence.TransactionHandler, id int) error
	RestoreCapturePage(ctx context.Context, tx persistence.TransactionHandler, id int) error
}
