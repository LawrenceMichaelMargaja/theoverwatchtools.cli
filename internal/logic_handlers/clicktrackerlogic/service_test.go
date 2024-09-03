package clicktrackerlogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/clicktrackerlogic/clicktrackerlogicfakes"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlconn"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistencefakes"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore/testhelper"
	"github.com/friendsofgo/errors"
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

type argsGetClickTrackers struct {
	User               mysqlmodel.User
	ClickTracker       mysqlmodel.ClickTracker
	ClickTrackerSet    mysqlmodel.ClickTrackerSet
	filter             *model.ClickTrackerFilters
	CreateClickTracker *model.CreateClickTracker
}

type testCaseGetClickTrackers struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsGetClickTrackers
	mutations       func(t *testing.T, db *sqlx.DB, args *argsGetClickTrackers)
	assertions      func(t *testing.T, organization *model.PaginatedClickTrackers, err error)
}

func getGetClickTrackerTestCases() []testCaseGetClickTrackers {
	return []testCaseGetClickTrackers{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsGetClickTrackers{
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					Name:          "Younes",
					CreatedBy:     null.IntFrom(2),
					LastUpdatedBy: null.IntFrom(2),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					Name:              "Demby",
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					ClickTrackerSetID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetClickTrackers) {
				err := args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
			assertions: func(t *testing.T, paginated *model.PaginatedClickTrackers, err error) {
				require.NoError(t, err, "unexpected get organizations error")
				require.NotNil(t, paginated, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click trackers")
				assert.NotEqual(t, 0, paginated.Pagination.RowCount, "unexpected total count")
			},
		},
		{
			name:            "fail-get-click-trackers",
			getDependencies: getConcreteDependencies,
			args: &argsGetClickTrackers{
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					Name:          "Younes",
					CreatedBy:     null.IntFrom(2),
					LastUpdatedBy: null.IntFrom(2),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					Name:              "Demby",
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					ClickTrackerSetID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetClickTrackers) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.ClickTracker)
			},
			assertions: func(t *testing.T, paginated *model.PaginatedClickTrackers, err error) {
				require.Error(t, err, "unexpected get click trackers error")
				require.Contains(t, err.Error(), "get click trackers:")
			},
		},
		{
			name: "fail-mock-get-db",
			getDependencies: func(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
				cleanup := func(ignoreErrors ...bool) {

				}

				mockTxProvider := persistencefakes.FakeTransactionProvider{}
				mockTxProvider.DbReturns(nil, errors.New(mockDbReturnsErr))

				return &dependencies{
					Persistor:  &clicktrackerlogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			args: &argsGetClickTrackers{
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					Name:          "Younes",
					CreatedBy:     null.IntFrom(2),
					LastUpdatedBy: null.IntFrom(2),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					Name:              "Demby",
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					ClickTrackerSetID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetClickTrackers) {
			},
			assertions: func(t *testing.T, paginated *model.PaginatedClickTrackers, err error) {
				require.Error(t, err, "unexpected get organizations error")
				require.Contains(t, err.Error(), "get db:")
			},
		},
	}
}

func TestService_GetClickTrackers(t *testing.T) {
	for _, tt := range getGetClickTrackerTestCases() {
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

			paginatedClickTrackers, err := svc.ListClickTrackers(context.Background(), tt.args.filter)
			tt.assertions(t, paginatedClickTrackers, err)
		})
	}
}

type argsCreateClickTracker struct {
	ClickTracker       mysqlmodel.ClickTracker
	ClickTrackerSet    mysqlmodel.ClickTrackerSet
	CreateClickTracker *model.CreateClickTracker
}

type testCaseCreateClickTracker struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsCreateClickTracker
	mutations       func(t *testing.T, db *sqlx.DB, args *argsCreateClickTracker)
	assertions      func(t *testing.T, clickTracker *model.ClickTracker, err error)
}

func getTestCasesCreateClickTracker() []testCaseCreateClickTracker {
	return []testCaseCreateClickTracker{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsCreateClickTracker{
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            3,
					Name:          "Mohamed",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                1,
					Name:              "Demby",
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					ClickTrackerSetID: 3,
				},
				CreateClickTracker: &model.CreateClickTracker{
					Name:              "Younes",
					UserId:            1,
					ClickTrackerSetId: 3,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsCreateClickTracker) {
				err := args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
			assertions: func(t *testing.T, clickTracker *model.ClickTracker, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, clickTracker, "unexpected nil click tracker")

				require.NotEqual(t, 0, clickTracker.Id, "unexpected nil click tracker")
				require.NotEmpty(t, clickTracker.Name, "unexpected empty click tracker page name")
				require.NotEqual(t, 0, clickTracker.CreatedBy, "unexpected empty Created_by")
			},
		},
	}
}

func TestService_CreateClickTracker(t *testing.T) {
	for _, tt := range getTestCasesCreateClickTracker() {
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

			capturePage, err := svc.AddClickTracker(context.Background(), tt.args.CreateClickTracker)
			tt.assertions(t, capturePage, err)
		})
	}
}
