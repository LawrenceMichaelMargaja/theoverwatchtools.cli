package mysqlstore

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strings"
	"testing"
	"time"
)

type argsGetClickTrackers struct {
	User            mysqlmodel.User
	ClickTrackerSet mysqlmodel.ClickTrackerSet
	ClickTracker    mysqlmodel.ClickTracker
	ClickTrackerTwo mysqlmodel.ClickTracker
}

type testCaseGetClickTrackers struct {
	name       string
	filter     *model.ClickTrackerFilters
	args       *argsGetClickTrackers
	mutations  func(t *testing.T, db *sqlx.DB, args *argsGetClickTrackers) []int
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error, ids []int)
}

func getTestCasesGetClickTrackers() []testCaseGetClickTrackers {
	return []testCaseGetClickTrackers{
		{
			name: "success-filter-ids-in",
			filter: &model.ClickTrackerFilters{
				IdsIn: []int{4, 3},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click trackers")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.Equal(t, len(ids), len(paginated.ClickTrackers), "unexpected number of click trackers returned")

				for _, ct := range paginated.ClickTrackers {
					assert.Contains(t, ids, ct.Id, "returned click tracker ID not found in expected IDs")
				}

				assert.Equal(t, len(ids), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
			args: &argsGetClickTrackers{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            3,
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                4,
					Name:              "TEST",
					ClickTrackerSetID: 3,
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					CreatedAt:         time.Now(),
					LastUpdatedAt:     null.TimeFrom(time.Now()),
					IsActive:          true,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetClickTrackers) []int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				Ids := []int{args.User.ID}

				return Ids
			},
		},
		{
			name: "success-filter-names-in",
			filter: &model.ClickTrackerFilters{
				ClickTrackerNameIn: []string{"Lawrence", "Younes"},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click trackers")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.Equal(t, len(ids), len(paginated.ClickTrackers), "unexpected number of click trackers returned")

				for _, ct := range paginated.ClickTrackers {
					assert.Contains(t, ids, ct.Id, "returned click tracker ID not found in expected IDs")
				}

				assert.Equal(t, len(ids), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
			args: &argsGetClickTrackers{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            3,
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                4,
					Name:              "Lawrence",
					ClickTrackerSetID: 3,
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					CreatedAt:         time.Now(),
					LastUpdatedAt:     null.TimeFrom(time.Now()),
					IsActive:          true,
				},
				ClickTrackerTwo: mysqlmodel.ClickTracker{
					ID:                2,
					Name:              "Younes",
					ClickTrackerSetID: 3,
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					CreatedAt:         time.Now(),
					LastUpdatedAt:     null.TimeFrom(time.Now()),
					IsActive:          true,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetClickTrackers) []int {

				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTrackerTwo.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				Ids := []int{args.ClickTracker.ID, args.ClickTrackerTwo.ID}

				return Ids
			},
		},
		{
			name: "success-multiple-filters",
			filter: &model.ClickTrackerFilters{
				ClickTrackerNameIn: []string{"Demby"},
				IdsIn:              []int{1},
				CreatedBy:          null.IntFrom(4),
			},
			args: &argsGetClickTrackers{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            3,
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                1,
					Name:              "Demby",
					ClickTrackerSetID: 3,
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					CreatedAt:         time.Now(),
					LastUpdatedAt:     null.TimeFrom(time.Now()),
					IsActive:          true,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetClickTrackers) []int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting first click tracker")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting first click tracker")

				Ids := []int{args.ClickTracker.ID}

				return Ids
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click trackers")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.ClickTrackers) == 1, "unexpected greater than 1 click tracker")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 click tracker")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
		},
		{
			name: "empty-results",
			filter: &model.ClickTrackerFilters{
				ClickTrackerNameIn:   []string{"Saul Goodman"},
				IdsIn:                []int{1},
				ClickTrackerIsActive: null.BoolFrom(false),
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetClickTrackers) []int {
				return nil
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedClickTrackers, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.ClickTrackers, "unexpected nil click trackers")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.ClickTrackers) == 0, "unexpected greater than 1 click tracker")
				assert.True(t, paginated.Pagination.RowCount == 0, "unexpected count to be greater than 1 click tracker")

				modelhelpers.AssertNonEmptyClickTrackers(t, paginated.ClickTrackers)
			},
		},
	}
}

func Test_GetClickTrackers(t *testing.T) {
	for _, testCase := range getTestCasesGetClickTrackers() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			id := testCase.mutations(t, db, testCase.args)
			paginated, err := m.GetClickTrackers(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err, id)
		})
	}
}

type argsCreateClickTracker struct {
	User               mysqlmodel.User
	ClickTrackerSet    mysqlmodel.ClickTrackerSet
	CreateClickTracker *model.CreateClickTracker
}

type testCaseCreateClickTracker struct {
	name       string
	args       *argsCreateClickTracker
	assertions func(t *testing.T, db *sqlx.DB, clicktracker *model.ClickTracker, err error)
	mutations  func(t *testing.T, db *sqlx.DB, args *argsCreateClickTracker) *model.CreateClickTracker
}

func getCreateClickTrackerTestCases() []testCaseCreateClickTracker {
	return []testCaseCreateClickTracker{
		{
			name: "success",
			args: &argsCreateClickTracker{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					Name:          "lawrence",
					ID:            4,
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				CreateClickTracker: &model.CreateClickTracker{
					Name:              "younes",
					UserId:            4,
					ClickTrackerSetId: 4,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsCreateClickTracker) *model.CreateClickTracker {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the click tracker set db")

				return args.CreateClickTracker
			},
			assertions: func(t *testing.T, db *sqlx.DB, clickTracker *model.ClickTracker, err error) {
				assert.NotNil(t, clickTracker, "unexpected nil click tracker")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyClickTrackers(t, []model.ClickTracker{*clickTracker})
			},
		},
		{
			name: "fail-name-too-long",
			args: &argsCreateClickTracker{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					Name:          "lawrence",
					ID:            4,
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				CreateClickTracker: &model.CreateClickTracker{
					Name:              "DembyYounesLawrence",
					UserId:            4,
					ClickTrackerSetId: 4,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsCreateClickTracker) *model.CreateClickTracker {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the click tracker set db")

				repeatedName := strings.Repeat(args.CreateClickTracker.Name, 100)

				args.CreateClickTracker.Name = repeatedName
				return args.CreateClickTracker
			},
			assertions: func(t *testing.T, db *sqlx.DB, clickTracker *model.ClickTracker, err error) {
				assert.Nil(t, clickTracker, "unexpected non-nil click tracker")
				assert.Error(t, err, "expected an error due to name exceeding length limit")
			},
		},
	}
}

func Test_CreateClickTracker(t *testing.T) {
	for _, testCase := range getCreateClickTrackerTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			res := testCase.mutations(t, db, testCase.args)

			createClickTracker, err := m.AddClickTracker(testCtx, txHandlerDb, res)
			testCase.assertions(t, db, createClickTracker, err)
		})
	}
}

type argsDeleteClickTracker struct {
	User            mysqlmodel.User
	ClickTrackerSet mysqlmodel.ClickTrackerSet
	ClickTracker    mysqlmodel.ClickTracker
}

type testCaseDeleteClickTrackers struct {
	name       string
	args       *argsDeleteClickTracker
	assertions func(t *testing.T, db *sqlx.DB, id int)
	mutations  func(t *testing.T, db *sqlx.DB, args *argsDeleteClickTracker)
}

func getTestCasesDeleteClickTrackers() []testCaseDeleteClickTrackers {
	return []testCaseDeleteClickTrackers{
		{
			name: "success-filter-ids-in",
			args: &argsDeleteClickTracker{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            4,
					Name:          "lawrence",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                1,
					Name:              "Lawrence",
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					ClickTrackerSetID: 4,
					IsActive:          true,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteClickTracker) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the click tracker set db")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting the click tracker in the click tracker table.")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				entry, err := mysqlmodel.FindClickTracker(context.TODO(), db, id)
				require.Nil(t, err, "unexpected non-nil error")
				require.NoError(t, err, "unexpected error fetching the click tracker")

				assert.Equal(t, false, entry.IsActive)
			},
		},
		{
			name: "failure-non-existent-organization",
			args: &argsDeleteClickTracker{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            4,
					Name:          "lawrence",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                1,
					Name:              "Lawrence",
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					ClickTrackerSetID: 4,
					IsActive:          true,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				_, err := mysqlmodel.FindClickTracker(context.TODO(), db, id)
				require.Error(t, err, "expected an error for non-existent click tracker")
				require.Contains(t, err.Error(), "no rows", "error should indicate that the organization was not found")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteClickTracker) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
		},
	}
}

func Test_DeleteClickTracker(t *testing.T) {
	for _, testCase := range getTestCasesDeleteClickTrackers() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			testCase.mutations(t, db, testCase.args)

			err = m.DeleteClickTracker(testCtx, txHandlerDb, testCase.args.ClickTracker.ID)
			require.NoError(t, err, "error deleting the click tracker.")
			testCase.assertions(t, db, testCase.args.ClickTracker.ID)
		})
	}
}

type argsUpdateClickTracker struct {
	User                   mysqlmodel.User
	ClickTrackerSet        mysqlmodel.ClickTrackerSet
	ClickTracker           mysqlmodel.ClickTracker
	UpdateClickTrackerData *model.UpdateClickTracker
}

type testCaseUpdateClickTrackers struct {
	name       string
	args       *argsUpdateClickTracker
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB, args *argsUpdateClickTracker)
}

func getTestCasesUpdateClickTrackers() []testCaseUpdateClickTrackers {
	return []testCaseUpdateClickTrackers{
		{
			name: "success-filter-ids-in",
			args: &argsUpdateClickTracker{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            4,
					Name:          "lawrence",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                1,
					Name:              "Lawrence",
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					ClickTrackerSetID: 4,
					IsActive:          true,
				},
				UpdateClickTrackerData: &model.UpdateClickTracker{
					Id: 1,
					Name: null.String{
						String: "Demby",
						Valid:  true,
					},
					ClickTrackerSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateClickTracker) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the click tracker set db")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting the click tracker in the click tracker table.")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindClickTracker(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the click tracker")

				assert.Equal(t, true, entry.IsActive)
			},
		},
		{
			name: "failure-non-existent-organization",
			args: &argsUpdateClickTracker{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            4,
					Name:          "lawrence",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                1,
					Name:              "Lawrence",
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					ClickTrackerSetID: 4,
					IsActive:          true,
				},
				UpdateClickTrackerData: &model.UpdateClickTracker{
					Id: 1,
					Name: null.String{
						String: "Demby",
						Valid:  true,
					},
					ClickTrackerSetId: 1,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "expected error when updating a non-existent click tracker")
				assert.Contains(t, err.Error(), "expected exactly one entry for entity: 1")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateClickTracker) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
		},
	}
}

func Test_UpdateClickTracker(t *testing.T) {
	for _, testCase := range getTestCasesUpdateClickTrackers() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			testCase.mutations(t, db, testCase.args)

			_, err = m.UpdateClickTracker(testCtx, txHandlerDb, testCase.args.UpdateClickTrackerData)
			testCase.assertions(t, db, testCase.args.UpdateClickTrackerData.Id, err)
		})
	}
}

type argsRestoreClickTracker struct {
	User            mysqlmodel.User
	ClickTrackerSet mysqlmodel.ClickTrackerSet
	ClickTracker    mysqlmodel.ClickTracker
}

type testCaseRestoreClickTrackers struct {
	name       string
	args       *argsRestoreClickTracker
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB, args *argsRestoreClickTracker)
}

func getTestCasesRestoreClickTrackers() []testCaseRestoreClickTrackers {
	return []testCaseRestoreClickTrackers{
		{
			name: "success-filter-ids-in",
			args: &argsRestoreClickTracker{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            4,
					Name:          "lawrence",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                1,
					Name:              "Lawrence",
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					ClickTrackerSetID: 4,
					IsActive:          false,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreClickTracker) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.ClickTrackerSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the click tracker set db")

				err = args.ClickTracker.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting the click tracker in the click tracker table.")
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				clickTracker, err := mysqlmodel.FindClickTracker(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the click tracker")
				assert.Equal(t, true, clickTracker.IsActive)
			},
		},
		{
			name: "failure-non-existent-organization",
			args: &argsRestoreClickTracker{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				ClickTrackerSet: mysqlmodel.ClickTrackerSet{
					ID:            4,
					Name:          "lawrence",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				ClickTracker: mysqlmodel.ClickTracker{
					ID:                1,
					Name:              "Lawrence",
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					ClickTrackerSetID: 4,
					IsActive:          true,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				_, fetchErr := mysqlmodel.FindClickTracker(context.TODO(), db, id)
				assert.Error(t, fetchErr, "click tracker should not exist")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreClickTracker) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
		},
	}
}

func Test_RestoreClickTracker(t *testing.T) {
	for _, testCase := range getTestCasesRestoreClickTrackers() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			cfg := &Config{
				Logger:        testLogger,
				QueryTimeouts: testQueryTimeouts,
			}

			m, err := New(cfg)
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, m, "unexpected nil")

			txHandler, err := mysqltx.New(&mysqltx.Config{
				Logger:       testLogger,
				Db:           db,
				DatabaseName: cp.Database,
			})
			require.NoError(t, err, "unexpected error creating the tx handler")

			txHandlerDb, err := txHandler.Db(testCtx)
			require.NoError(t, err, "unexpected error fetching the db from the tx handler")
			require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")

			testCase.mutations(t, db, testCase.args)

			err = m.RestoreClickTracker(testCtx, txHandlerDb, testCase.args.ClickTracker.ID)
			require.NoError(t, err, "error deleting the click tracker.")
			testCase.assertions(t, db, testCase.args.ClickTracker.ID, err)
		})
	}
}
