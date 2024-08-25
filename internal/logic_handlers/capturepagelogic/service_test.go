package capturepagelogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"testing"
	"time"
)

var (
	mockTimeout      = 5 * time.Second
	mockLogger       = logger.New(context.TODO())
	mockDbReturnsErr = "error getting db"
)

type dependencies struct {
	Persistor  persistor
	Logger     *logrus.Entry
	TxProvider persistence.TransactionProvider
	Db         *sqlx.DB
	Cleanup    func(ignoreErrors ...bool)
}

func getConcreteDependencies(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)

	store, err := mysqlstore.New(&mysqlstore.Config{
		Logger: mockLogger,
		QueryTimeouts: &persistence.QueryTimeouts{
			Query: mockTimeout,
			Exec:  mockTimeout,
		},
	})
	require.NoError(t, err, "unexpected new mysqlstore error")

	tx, err := mysqltx.New(&mysqltx.Config{
		Logger:       mockLogger,
		Db:           db,
		DatabaseName: cp.Database,
	})
	require.NoError(t, err, "unexpected new mysqltx error")

	prov, err := mysqlconn.New(&mysqlconn.Config{
		Logger:    mockLogger,
		TxHandler: tx,
	})
	require.NoError(t, err, "unexpected new mysqlconn error")

	return &dependencies{
		Persistor:  store,
		TxProvider: prov,
		Logger:     mockLogger,
		Cleanup:    cleanup,
		Db:         db,
	}, cleanup
}

type argsGetCapturePages struct {
	filter            *model.CapturePageFilters
	User              mysqlmodel.CapturePage
	CreateCapturePage *model.CreateCapturePage
	CapturePageSet    mysqlmodel.CapturePageSet
}

type testCaseGetCapturePages struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsGetCapturePages
	mutations       func(t *testing.T, db *sqlx.DB, args *argsGetCapturePages)
	assertions      func(t *testing.T, capturePage *model.PaginatedCapturePages, err error)
}

func getGetCapturePagesTestCases() []testCaseGetCapturePages {
	return []testCaseGetCapturePages{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsGetCapturePages{
				User: mysqlmodel.CapturePage{
					Name:             "Demby",
					CreatedBy:        null.IntFrom(1),
					LastUpdatedBy:    null.IntFrom(1),
					CapturePageSetID: 1,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:          "Younes",
					CreatedBy:     null.IntFrom(2),
					LastUpdatedBy: null.IntFrom(2),
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetCapturePages) {
				err := args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting into capture_page_set")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting into the capture_page db")

			},
			assertions: func(t *testing.T, paginated *model.PaginatedCapturePages, err error) {
				require.NoError(t, err, "unexpected get organizations error")
				require.NotNil(t, paginated, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				require.NotNil(t, paginated.CapturePages, "unexpected nil organizations")
				assert.NotEqual(t, 0, paginated.Pagination.RowCount, "unexpected total count")
			},
		},
	}
}

func TestService_GetCapturePages(t *testing.T) {
	for _, tt := range getGetCapturePagesTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			_dependencies, cleanup := tt.getDependencies(t)
			defer cleanup()

			svc, err := New(&Config{
				TxProvider: _dependencies.TxProvider,
				Logger:     _dependencies.Logger,
				Persistor:  _dependencies.Persistor,
			})
			require.NoError(t, err, "unexpected new error")

			tt.mutations(t, _dependencies.Db, tt.args)

			paginatedCapturePages, err := svc.ListCapturePages(context.Background(), tt.args.filter)
			tt.assertions(t, paginatedCapturePages, err)
		})
	}
}
