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
	"strconv"
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
	assertions      func(t *testing.T, clickTracker *model.PaginatedClickTrackers, err error)
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
				require.NoError(t, err, "unexpected get click tracker error")
				require.NotNil(t, paginated, "unexpected nil click trackers")
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
				require.Error(t, err, "unexpected get click trackers error")
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

			clickTracker, err := svc.AddClickTrackers(context.Background(), tt.args.CreateClickTracker)
			tt.assertions(t, clickTracker, err)
		})
	}
}

type argsUpdateClickTracker struct {
	ClickTracker       mysqlmodel.ClickTracker
	ClickTrackerSet    mysqlmodel.ClickTrackerSet
	UpdateClickTracker *model.UpdateClickTracker
}

type testCaseUpdateClickTracker struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsUpdateClickTracker
	mutations       func(t *testing.T, db *sqlx.DB, args *argsUpdateClickTracker)
	assertions      func(t *testing.T, params *model.UpdateClickTracker, clickTracker *model.ClickTracker, err error)
}

func getTestCasesUpdateClickTracker() []testCaseUpdateClickTracker {
	return []testCaseUpdateClickTracker{
		{
			name: "success",
			args: &argsUpdateClickTracker{
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
				UpdateClickTracker: &model.UpdateClickTracker{
					Id:     1,
					Name:   null.StringFrom("Younes"),
					UserId: null.IntFrom(1),
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateClickTracker) {
				err := args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, params *model.UpdateClickTracker, clickTracker *model.ClickTracker, err error) {
				createdByConvToInt, _ := strconv.Atoi(clickTracker.CreatedBy)
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, clickTracker, "unexpected nil click tracker")
				assert.Equal(t, params.Id, clickTracker.Id, "expected id to be equal")
				assert.Equal(t, params.UserId.Int, createdByConvToInt, "expected created_by id to be equal")
				assert.Equal(t, params.Name.String, clickTracker.Name, "expected name to be equal")
			},
		},
		{
			name: "success-name-only",
			args: &argsUpdateClickTracker{
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
				UpdateClickTracker: &model.UpdateClickTracker{
					Id:     1,
					Name:   null.StringFrom("Younes"),
					UserId: null.IntFrom(1),
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateClickTracker) {
				err := args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, params *model.UpdateClickTracker, clickTracker *model.ClickTracker, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, clickTracker, "unexpected nil click tracker")
				assert.Equal(t, params.Name.String, clickTracker.Name, "expected name to be equal")
			},
		},
		{
			name: "fail-mock",
			args: &argsUpdateClickTracker{
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
				UpdateClickTracker: &model.UpdateClickTracker{
					Id:     1,
					Name:   null.StringFrom("Younes"),
					UserId: null.IntFrom(1),
				},
			},
			getDependencies: func(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
				cleanup := func(ignoreErrors ...bool) {}

				mockTxProvider := persistencefakes.FakeTransactionProvider{}
				mockTxProvider.TxReturns(nil, errors.New(mockDbReturnsErr))

				return &dependencies{
					Persistor:  &clicktrackerlogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateClickTracker) {
			},
			assertions: func(t *testing.T, params *model.UpdateClickTracker, clickTracker *model.ClickTracker, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, clickTracker, "unexpected nil click tracker")
			},
		},
		{
			name: "fail-internal-server-error",
			args: &argsUpdateClickTracker{
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
				UpdateClickTracker: &model.UpdateClickTracker{
					Id:     1,
					Name:   null.StringFrom("Younes"),
					UserId: null.IntFrom(1),
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateClickTracker) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.ClickTracker)
			},
			assertions: func(t *testing.T, params *model.UpdateClickTracker, clickTracker *model.ClickTracker, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, clickTracker, "unexpected nil click tracker")
			},
		},
	}
}

func TestService_UpdateClickTracker(t *testing.T) {
	for _, tt := range getTestCasesUpdateClickTracker() {
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

			clickTracker, err := svc.UpdateClickTracker(context.Background(), tt.args.UpdateClickTracker)
			tt.assertions(t, tt.args.UpdateClickTracker, clickTracker, err)
		})
	}
}

type argsDeleteClickTracker struct {
	User               mysqlmodel.User
	ClickTrackerSet    mysqlmodel.ClickTrackerSet
	DeleteClickTracker *model.DeleteClickTracker
}

type testCaseDeleteClickTracker struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsDeleteClickTracker
	mutations       func(t *testing.T, db *sqlx.DB, args *argsDeleteClickTracker)
	assertions      func(t *testing.T, db *sqlx.DB, id int)
}

func getTestCasesDeleteClickTracker() []testCaseDeleteClickTracker {
	return []testCaseDeleteClickTracker{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsDeleteClickTracker{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            1,
					Name:          "Mohamed",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				DeleteClickTracker: &model.DeleteClickTracker{
					ID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteClickTracker) {
				err := args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				returnClickTracker, err := mysqlmodel.FindClickTracker(context.TODO(), db, id)
				require.Nil(t, returnClickTracker, "expected to be nil")
				require.Error(t, err, "error fetching click tracker from db")
			},
		},
		{
			name:            "fail-click-tracker-not-found",
			getDependencies: getConcreteDependencies,
			args: &argsDeleteClickTracker{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            4,
					Name:          "Mohamed",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				DeleteClickTracker: &model.DeleteClickTracker{
					ID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteClickTracker) {
				err := args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				returnClickTracker, err := mysqlmodel.FindClickTracker(context.TODO(), db, id)
				require.Nil(t, returnClickTracker, "expected to be nil")
				require.Error(t, err, "expected an error when deleting a non-existent click tracker")
			},
		},
	}
}

func TestService_DeleteClickTracker(t *testing.T) {
	for _, testCase := range getTestCasesDeleteClickTracker() {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			_dependencies, cleanup := testCase.getDependencies(t)
			defer cleanup()

			svc, err := New(&Config{
				TxProvider: _dependencies.TxProvider,
				Logger:     _dependencies.Logger,
				Persistor:  _dependencies.Persistor,
			})
			require.NoError(t, err, "unexpected new error")

			testCase.mutations(t, _dependencies.Db, testCase.args)
			require.NoError(t, err, "unexpected new service error")

			err = svc.DeleteClickTracker(context.Background(), testCase.args.DeleteClickTracker)
			require.NoError(t, err, "unexpected error deleting click tracker.")
			testCase.assertions(t, db, testCase.args.DeleteClickTracker.ID)
		})
	}
}

type argsRestoreClickTracker struct {
	User                mysqlmodel.User
	ClickTracker        mysqlmodel.ClickTracker
	ClickTrackerSet     mysqlmodel.ClickTrackerSet
	RestoreClickTracker *model.RestoreClickTracker
}

type testCaseRestoreClickTracker struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsRestoreClickTracker
	mutations       func(t *testing.T, db *sqlx.DB, args *argsRestoreClickTracker)
	assertions      func(t *testing.T, db *sqlx.DB, id int, err error)
}

func getTestCasesRestoreClickTracker() []testCaseRestoreClickTracker {
	return []testCaseRestoreClickTracker{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsRestoreClickTracker{
				User: mysqlmodel.User{
					ID:                6,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                3,
					Name:              "Mohamed",
					ClickTrackerSetID: 1,
					CreatedBy:         null.IntFrom(2),
					LastUpdatedBy:     null.IntFrom(2),
					IsActive:          false,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            1,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				RestoreClickTracker: &model.RestoreClickTracker{
					ID: 3,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreClickTracker) {
				err := args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the click_tracker_set db")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the click_tracker db")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.NoError(t, err, "error restoring click tracker from db")

				returnClickTracker, err := mysqlmodel.FindClickTracker(context.TODO(), db, id)
				require.NoError(t, err, "error fetching click tracker from db")
				require.NotNil(t, returnClickTracker, "expected click tracker to be not nil")

				require.True(t, returnClickTracker.IsActive, "expected click tracker to be active")
				assert.Equal(t, returnClickTracker.IsActive, true)
			},
		},
		{
			name:            "fail-click-tracker-not-found",
			getDependencies: getConcreteDependencies,
			args: &argsRestoreClickTracker{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                3,
					Name:              "Mohamed",
					ClickTrackerSetID: 1,
					CreatedBy:         null.IntFrom(2),
					LastUpdatedBy:     null.IntFrom(2),
					IsActive:          false,
				},
				RestoreClickTracker: &model.RestoreClickTracker{
					ID: 32345,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreClickTracker) {

				err := args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the click_tracker_set db")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				returnClickTracker, fetchErr := mysqlmodel.FindClickTracker(context.TODO(), db, id)
				require.Error(t, fetchErr, "expected an error when fetching click tracker from db")
				require.Nil(t, returnClickTracker, "expected click tracker to be nil")
			},
		},
	}
}

func TestService_RestoreClickTracker(t *testing.T) {
	for _, testCase := range getTestCasesRestoreClickTracker() {
		t.Run(testCase.name, func(t *testing.T) {
			_dependencies, cleanup := testCase.getDependencies(t)
			defer cleanup()

			svc, err := New(&Config{
				TxProvider: _dependencies.TxProvider,
				Logger:     _dependencies.Logger,
				Persistor:  _dependencies.Persistor,
			})
			require.NoError(t, err, "unexpected new error")

			testCase.mutations(t, _dependencies.Db, testCase.args)
			require.NoError(t, err, "unexpected error in mutations")

			err = svc.RestoreClickTracker(context.Background(), testCase.args.RestoreClickTracker)
			testCase.assertions(t, _dependencies.Db, testCase.args.RestoreClickTracker.ID, err)
		})
	}
}
