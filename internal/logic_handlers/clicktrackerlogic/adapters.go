package clicktrackerlogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . persistor
type persistor interface {
	GetClickTrackers(ctx context.Context, tx persistence.TransactionHandler, filters *model.ClickTrackerFilters) (*model.PaginatedClickTrackers, error)
	AddClickTracker(ctx context.Context, tx persistence.TransactionHandler, clickTracker *model.CreateClickTracker) (*model.ClickTracker, error)
	GetClickTrackerByName(ctx context.Context, tx persistence.TransactionHandler, name string) (*model.ClickTracker, error)
	UpdateClickTracker(ctx context.Context, tx persistence.TransactionHandler, params *model.UpdateClickTracker) (*model.ClickTracker, error)
	DeleteClickTracker(ctx context.Context, tx persistence.TransactionHandler, id int) error
	RestoreClickTracker(ctx context.Context, tx persistence.TransactionHandler, id int) error
}
