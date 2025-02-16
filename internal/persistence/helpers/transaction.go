package helpers

import "context"

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback() error
}
