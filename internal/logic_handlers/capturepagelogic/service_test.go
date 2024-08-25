package capturepagelogic

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/logic_handlers/capturepagelogic/capturepagelogicfakes"
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
				require.NoError(t, err, "unexpected get capture page error")
				require.NotNil(t, paginated, "unexpected nil capture pages")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				require.NotNil(t, paginated.CapturePages, "unexpected nil capture pages")
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

type argsUpdateCapturePage struct {
	User              mysqlmodel.CapturePage
	CapturePage       mysqlmodel.CapturePage
	CapturePageSet    mysqlmodel.CapturePageSet
	UpdateCapturePage *model.UpdateCapturePage
}

type testCaseUpdateCapturePage struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsUpdateCapturePage
	mutations       func(t *testing.T, db *sqlx.DB, args *argsUpdateCapturePage)
	assertions      func(t *testing.T, params *model.UpdateCapturePage, capturePage *model.CapturePage, err error)
}

func getTestCasesUpdateCapturePage() []testCaseUpdateCapturePage {
	return []testCaseUpdateCapturePage{
		{
			name: "success",
			args: &argsUpdateCapturePage{
				User: mysqlmodel.CapturePage{
					ID:               1,
					Name:             "Demby",
					CreatedBy:        null.IntFrom(1),
					LastUpdatedBy:    null.IntFrom(1),
					CapturePageSetID: 1,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:          "Mohamed",
					CreatedBy:     null.IntFrom(3),
					LastUpdatedBy: null.IntFrom(3),
				},
				UpdateCapturePage: &model.UpdateCapturePage{
					Id:     1,
					Name:   null.StringFrom("Younes"),
					UserId: null.IntFrom(1),
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateCapturePage) {
				err := args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, params *model.UpdateCapturePage, capturePage *model.CapturePage, err error) {
				createdByConvToInt, _ := strconv.Atoi(capturePage.CreatedBy)
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, capturePage, "unexpected nil capture page")
				assert.Equal(t, params.Id, capturePage.Id, "expected id to be equal")
				assert.Equal(t, params.UserId.Int, createdByConvToInt, "expected created_by id to be equal")
				assert.Equal(t, params.Name.String, capturePage.Name, "expected name to be equal")
			},
		},
		{
			name: "success-name-only",
			args: &argsUpdateCapturePage{
				User: mysqlmodel.CapturePage{
					ID:               1,
					Name:             "Demby",
					CreatedBy:        null.IntFrom(1),
					LastUpdatedBy:    null.IntFrom(1),
					CapturePageSetID: 1,
				},
				UpdateCapturePage: &model.UpdateCapturePage{
					Id:     1,
					Name:   null.StringFrom("Younes"),
					UserId: null.IntFrom(1),
				},
			},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateCapturePage) {
				err := args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, params *model.UpdateCapturePage, capturePage *model.CapturePage, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, capturePage, "unexpected nil capture page")
				assert.Equal(t, params.Name.String, capturePage.Name, "expected name to be equal")
			},
		},
		{
			name: "fail-mock",
			args: &argsUpdateCapturePage{},
			getDependencies: func(t *testing.T) (*dependencies, func(ignoreErrors ...bool)) {
				cleanup := func(ignoreErrors ...bool) {}

				mockTxProvider := persistencefakes.FakeTransactionProvider{}
				mockTxProvider.TxReturns(nil, errors.New(mockDbReturnsErr))

				return &dependencies{
					Persistor:  &capturepagelogicfakes.FakePersistor{},
					TxProvider: &mockTxProvider,
					Logger:     mockLogger,
					Cleanup:    cleanup,
				}, cleanup
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateCapturePage) {},
			assertions: func(t *testing.T, params *model.UpdateCapturePage, capturePage *model.CapturePage, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, capturePage, "unexpected nil capture page")
			},
		},
		{
			name:            "fail-internal-server-error",
			args:            &argsUpdateCapturePage{},
			getDependencies: getConcreteDependencies,
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateCapturePage) {
				testhelper.DropTable(t, db, mysqlmodel.TableNames.CapturePage)
			},
			assertions: func(t *testing.T, params *model.UpdateCapturePage, capturePage *model.CapturePage, err error) {
				assert.Error(t, err, "unexpected error")
				assert.Nil(t, capturePage, "unexpected nil capture page")
			},
		},
	}
}

func TestService_UpdateCapturePage(t *testing.T) {
	for _, tt := range getTestCasesUpdateCapturePage() {
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

			capturePage, err := svc.UpdateCapturePage(context.Background(), tt.args.UpdateCapturePage)
			tt.assertions(t, tt.args.UpdateCapturePage, capturePage, err)
		})
	}
}

type argsCreateCapturePage struct {
	User              mysqlmodel.CapturePage
	CapturePageSet    mysqlmodel.CapturePageSet
	CreateCapturePage *model.CreateCapturePage
}

type testCaseCreateCapturePage struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsCreateCapturePage
	mutations       func(t *testing.T, db *sqlx.DB, args *argsCreateCapturePage)
	assertions      func(t *testing.T, capturePage *model.CapturePage, err error)
}

func getTestCasesCreateCapturePage() []testCaseCreateCapturePage {
	return []testCaseCreateCapturePage{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsCreateCapturePage{
				User: mysqlmodel.CapturePage{
					ID:               1,
					Name:             "Demby",
					CreatedBy:        null.IntFrom(1),
					LastUpdatedBy:    null.IntFrom(1),
					CapturePageSetID: 1,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					Name:          "Mohamed",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				CreateCapturePage: &model.CreateCapturePage{
					Name:             "Younes",
					UserId:           1,
					CapturePageSetId: 1,
				},
			},
			assertions: func(t *testing.T, capturePage *model.CapturePage, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, capturePage, "unexpected nil capture page")

				require.NotEqual(t, 0, capturePage.Id, "unexpected nil capture page")
				require.NotEmpty(t, capturePage.Name, "unexpected empty capture page name")
				require.NotEqual(t, 0, capturePage.CreatedBy, "unexpected empty Created_by")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsCreateCapturePage) {
				err := args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
		},
	}
}

func TestService_CreateCapturePage(t *testing.T) {
	for _, tt := range getTestCasesCreateCapturePage() {
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

			capturePage, err := svc.AddCapturePage(context.Background(), tt.args.CreateCapturePage)
			tt.assertions(t, capturePage, err)
		})
	}
}

type argsDeleteCapturePage struct {
	User              mysqlmodel.User
	CapturePageSet    mysqlmodel.CapturePageSet
	DeleteCapturePage *model.DeleteCapturePage
}

type testCaseDeleteCapturePage struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsDeleteCapturePage
	mutations       func(t *testing.T, db *sqlx.DB, args *argsDeleteCapturePage)
	assertions      func(t *testing.T, db *sqlx.DB, id int)
}

func getTestCasesDeleteCapturePage() []testCaseDeleteCapturePage {
	return []testCaseDeleteCapturePage{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsDeleteCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					ID:            1,
					Name:          "Mohamed",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				DeleteCapturePage: &model.DeleteCapturePage{
					ID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteCapturePage) {
				err := args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				returnCapturePage, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.Nil(t, returnCapturePage, "expected to be nil")
				require.Error(t, err, "error fetching capture page from db")
			},
		},
		{
			name:            "fail-capture-page-not-found",
			getDependencies: getConcreteDependencies,
			args: &argsDeleteCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					ID:            4,
					Name:          "Mohamed",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				DeleteCapturePage: &model.DeleteCapturePage{
					ID: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteCapturePage) {
				err := args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				returnCapturePage, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.Nil(t, returnCapturePage, "expected to be nil")
				require.Error(t, err, "expected an error when deleting a non-existent capture page")
			},
		},
	}
}

func TestService_DeleteCapturePage(t *testing.T) {
	for _, testCase := range getTestCasesDeleteCapturePage() {
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

			err = svc.DeleteCapturePage(context.Background(), testCase.args.DeleteCapturePage)
			require.NoError(t, err, "unexpected error deleting capture page.")
			testCase.assertions(t, db, testCase.args.DeleteCapturePage.ID)
		})
	}
}

type argsRestoreCapturePage struct {
	User               mysqlmodel.User
	CapturePage        mysqlmodel.CapturePage
	CapturePageSet     mysqlmodel.CapturePageSet
	RestoreCapturePage *model.RestoreCapturePage
}

type testCaseRestoreCapturePage struct {
	name            string
	getDependencies func(t *testing.T) (*dependencies, func(ignoreErrors ...bool))
	args            *argsRestoreCapturePage
	mutations       func(t *testing.T, db *sqlx.DB, args *argsRestoreCapturePage)
	assertions      func(t *testing.T, db *sqlx.DB, id int, err error)
}

func getTestCasesRestoreCapturePage() []testCaseRestoreCapturePage {
	return []testCaseRestoreCapturePage{
		{
			name:            "success",
			getDependencies: getConcreteDependencies,
			args: &argsRestoreCapturePage{
				User: mysqlmodel.User{
					ID:                6,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				CapturePage: mysqlmodel.CapturePage{
					ID:               3,
					Name:             "Mohamed",
					CapturePageSetID: 1,
					CreatedBy:        null.IntFrom(2),
					LastUpdatedBy:    null.IntFrom(2),
					IsActive:         false,
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					ID:            1,
					Name:          "Demby",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				RestoreCapturePage: &model.RestoreCapturePage{
					ID: 3,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreCapturePage) {
				err := args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the capture_page_set db")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.CapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the capture_page db")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.NoError(t, err, "error restoring capture page from db")

				returnCapturePage, err := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.NoError(t, err, "error fetching capture page from db")
				require.NotNil(t, returnCapturePage, "expected capture page to be not nil")

				require.True(t, returnCapturePage.IsActive, "expected capture page to be active")
				assert.Equal(t, returnCapturePage.IsActive, true)
			},
		},
		{
			name:            "fail-capture-page-not-found",
			getDependencies: getConcreteDependencies,
			args: &argsRestoreCapturePage{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				CapturePage: mysqlmodel.CapturePage{
					ID:               3,
					Name:             "Mohamed",
					CapturePageSetID: 1,
					CreatedBy:        null.IntFrom(2),
					LastUpdatedBy:    null.IntFrom(2),
					IsActive:         false,
				},
				RestoreCapturePage: &model.RestoreCapturePage{
					ID: 32345,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreCapturePage) {
				err := args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the capture_page_set db")

				err = args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.CapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				returnCapturePage, fetchErr := mysqlmodel.FindCapturePage(context.TODO(), db, id)
				require.Error(t, fetchErr, "expected an error when fetching capture page from db")
				require.Nil(t, returnCapturePage, "expected capture page to be nil")
			},
		},
	}
}

func TestService_RestoreCapturePage(t *testing.T) {
	for _, testCase := range getTestCasesRestoreCapturePage() {
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

			err = svc.RestoreCapturePage(context.Background(), testCase.args.RestoreCapturePage)
			testCase.assertions(t, _dependencies.Db, testCase.args.RestoreCapturePage.ID, err)
		})
	}
}
