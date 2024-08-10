package mysqlstore

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqltx"
	"github.com/dembygenesis/local.tools/internal/persistence/persistors/mysqlstore/testhelper"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strings"
	"testing"
	"time"
)

var (
	testCtx           = context.TODO()
	testLogger        = logger.New(context.TODO())
	testQueryTimeouts = &persistence.QueryTimeouts{
		Query: 10 * time.Second,
		Exec:  10 * time.Second,
	}
)

func TestNew(t *testing.T) {
	cfg := &Config{
		Logger:        testLogger,
		QueryTimeouts: testQueryTimeouts,
	}

	m, err := New(cfg)
	require.NoError(t, err, "unexpected error")
	require.NotNil(t, m, "unexpected nil")
}

type testCaseGetCategories struct {
	name       string
	filter     *model.CategoryFilters
	mutations  func(t *testing.T, db *sqlx.DB)
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error)
}

func getTestCasesGetCategories() []testCaseGetCategories {
	return []testCaseGetCategories{
		{
			name: "success-filter-ids-in",
			filter: &model.CategoryFilters{
				IdsIn: []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) == 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "success-filter-names-in",
			filter: &model.CategoryFilters{
				CategoryNameIn: []string{"Regular User", "Admin"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) == 2, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 2, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "success-category-type-id-in",
			filter: &model.CategoryFilters{
				CategoryTypeIdIn: []int{1},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) > 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "success-category-type-name-in",
			filter: &model.CategoryFilters{
				CategoryTypeNameIn: []string{"User Types"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) > 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount > 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "success-multiple-filters",
			filter: &model.CategoryFilters{
				CategoryTypeNameIn: []string{"User Types"},
				CategoryTypeIdIn:   []int{1},
				CategoryNameIn:     []string{"Super Admin"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) == 1, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
		{
			name: "empty-results",
			filter: &model.CategoryFilters{
				CategoryTypeNameIn: []string{"Saul Goodman"},
				CategoryTypeIdIn:   []int{1},
				CategoryNameIn:     []string{"Super Admin"},
			},
			mutations: func(t *testing.T, db *sqlx.DB) {

			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedCategories, err error) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Categories, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Categories) == 0, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 0, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyCategories(t, paginated.Categories)
			},
		},
	}
}

func Test_GetCategories(t *testing.T) {
	for _, testCase := range getTestCasesGetCategories() {
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
			testCase.mutations(t, db)
			paginated, err := m.GetCategories(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err)
		})
	}
}

func TestNew_Fail(t *testing.T) {
	cfg := &Config{
		Logger:        nil,
		QueryTimeouts: testQueryTimeouts,
	}

	m, err := New(cfg)
	require.Error(t, err, "unexpected nil error")
	require.Nil(t, m, "unexpected not nil")
}

func Test_ReadCategories(t *testing.T) {
	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
	defer cleanup()

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

	paginatedCategories, err := m.GetCategories(testCtx, txHandlerDb, nil)
	require.NoError(t, err, "unexpected error fetching the categories from the database")
	require.NotNil(t, txHandlerDb, "unexpected nil categories")
	require.True(t, len(paginatedCategories.Categories) > 0, "unexpected empty categories")
}

type testCaseUpdateCategories struct {
	name       string
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB) (updateData model.UpdateCategory)
}

func getUpdateCategoriesTestCases() []testCaseUpdateCategories {
	return []testCaseUpdateCategories{
		{
			name: "success",
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCategory(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, 1, entry.CategoryTypeRefID)
			},
			mutations: func(t *testing.T, db *sqlx.DB) (updateData model.UpdateCategory) {

				updateCategory := mysqlmodel.Category{
					ID:                4,
					CategoryTypeRefID: 1,
					Name:              "OLD NAME",
				}

				err := updateCategory.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				updateCategoryData := model.UpdateCategory{
					Id:                4,
					CategoryTypeRefId: null.IntFrom(1),
					Name:              null.StringFrom("Category A new"),
				}

				return updateCategoryData
			},
		},
		{
			name: "fail",
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "expected error when updating a non-existent organization")
				assert.Contains(t, err.Error(), "expected exactly one entry for entity: 124124")
			},
			mutations: func(t *testing.T, db *sqlx.DB) (updateData model.UpdateCategory) {

				updateCategory := mysqlmodel.Category{
					ID:                4,
					CategoryTypeRefID: 1,
					Name:              "OLD NAME",
				}

				err := updateCategory.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				updateCategoryData := model.UpdateCategory{
					Id:                124124,
					CategoryTypeRefId: null.IntFrom(1),
					Name:              null.StringFrom("Category A new"),
				}

				return updateCategoryData
			},
		},
	}
}

func Test_UpdateCategories(t *testing.T) {
	for _, testCase := range getUpdateCategoriesTestCases() {
		db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
		defer cleanup()

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

		updateData := testCase.mutations(t, db)

		_, err = m.UpdateCategory(testCtx, txHandlerDb, &updateData)
		testCase.assertions(t, db, updateData.Id, err)
	}
}

//func Test_UpdateCategories_Success(t *testing.T) {
//	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
//	defer cleanup()
//
//	cfg := &Config{
//		Logger:        testLogger,
//		QueryTimeouts: testQueryTimeouts,
//	}
//
//	m, err := New(cfg)
//	require.NoError(t, err, "unexpected error")
//	require.NotNil(t, m, "unexpected nil")
//
//	txHandler, err := mysqltx.New(&mysqltx.Config{
//		Logger:       testLogger,
//		Db:           db,
//		DatabaseName: cp.Database,
//	})
//	require.NoError(t, err, "unexpected error creating the tx handler")
//
//	txHandlerDb, err := txHandler.Db(testCtx)
//	require.NoError(t, err, "unexpected error fetching the db from the tx handler")
//	require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")
//
//	paginatedCategories, err := m.GetCategories(testCtx, txHandlerDb, nil)
//	require.NoError(t, err, "unexpected error fetching the categories from the database")
//	require.NotNil(t, txHandlerDb, "unexpected nil categories")
//	require.True(t, len(paginatedCategories.Categories) > 0, "unexpected empty categories")
//
//	cat, err := m.UpdateCategory(testCtx, txHandlerDb, &updateCategory)
//	require.NoError(t, err, "unexpected error updating a conflicting category from the database")
//	assert.Equal(t, paginatedCategories.Categories[0].Name+" new", cat.Name)
//}

//func Test_UpdateCategories_Fail(t *testing.T) {
//	db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
//	defer cleanup()
//
//	cfg := &Config{
//		Logger:        testLogger,
//		QueryTimeouts: testQueryTimeouts,
//	}
//
//	m, err := New(cfg)
//	require.NoError(t, err, "unexpected error")
//	require.NotNil(t, m, "unexpected nil")
//
//	txHandler, err := mysqltx.New(&mysqltx.Config{
//		Logger:       testLogger,
//		Db:           db,
//		DatabaseName: cp.Database,
//	})
//	require.NoError(t, err, "unexpected error creating the tx handler")
//
//	txHandlerDb, err := txHandler.Db(testCtx)
//	require.NoError(t, err, "unexpected error fetching the db from the tx handler")
//	require.NotNil(t, txHandlerDb, "unexpected nil tx handler db")
//
//	paginatedCategories, err := m.GetCategories(testCtx, txHandlerDb, nil)
//	require.NoError(t, err, "unexpected error fetching the categories from the database")
//	require.NotNil(t, txHandlerDb, "unexpected nil categories")
//	require.True(t, len(paginatedCategories.Categories) > 0, "unexpected empty categories")
//
//	updateCategory := model.UpdateCategory{
//		Id: paginatedCategories.Categories[1].Id,
//		CategoryTypeRefId: null.Int{
//			Int:   paginatedCategories.Categories[0].CategoryTypeRefId,
//			Valid: true,
//		},
//		Name: null.String{
//			String: paginatedCategories.Categories[0].Name,
//			Valid:  true,
//		},
//	}
//
//	cat, err := m.UpdateCategory(testCtx, txHandlerDb, &updateCategory)
//	require.Error(t, err, "unexpected nil error fetching a conflicting category from the database")
//	assert.Contains(t, err.Error(), "Duplicate entry")
//	assert.Nil(t, cat, "unexpected non nil entry")
//}

type testCaseRestoreCategory struct {
	name       string
	id         int
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB)
}

func getRestoreCategoryTestCases() []testCaseRestoreCategory {
	return []testCaseRestoreCategory{
		{
			name: "success",
			id:   1,
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCategory(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the category")

				assert.Equal(t, 1, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Category{
					CategoryTypeRefID: 1,
					Name:              "test",
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")

				entry.IsActive = 0
				_, err = entry.Update(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected update error")

				err = entry.Reload(context.TODO(), db)
				assert.NoError(t, err, "unexpected reload error")

				assert.Equal(t, 0, entry.IsActive)
			},
		},
		{
			name: "fail-missing-entry-to-update",
			id:   1,
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "unexpected non-nil error")
				assert.Contains(t, err.Error(), "restore:")
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				testhelper.DropTable(t, db, "category")
			},
		},
	}
}

func Test_RestoreCategory(t *testing.T) {
	for _, testCase := range getRestoreCategoryTestCases() {
		db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
		t.Run(testCase.name, func(t *testing.T) {
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			defer cleanup()
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

			testCase.mutations(t, db)
			err = m.RestoreCategory(testCtx, txHandlerDb, testCase.id)
			testCase.assertions(t, db, testCase.id, err)
		})
	}
}

type testCaseCreateCategory struct {
	name          string
	categoryName  string
	categoryRefId int
	assertions    func(t *testing.T, db *sqlx.DB, category *model.Category, err error)
}

func getAddCategoryTestCases() []testCaseCreateCategory {
	return []testCaseCreateCategory{
		{
			name:          "success",
			categoryName:  "Example Category",
			categoryRefId: 1,
			assertions: func(t *testing.T, db *sqlx.DB, category *model.Category, err error) {
				assert.NotNil(t, category, "unexpected nil category")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyCategories(t, []model.Category{*category})
			},
		},
		{
			name:          "fail-name-exceeds-limit",
			categoryName:  strings.Repeat("a", 256),
			categoryRefId: 1,
			assertions: func(t *testing.T, db *sqlx.DB, category *model.Category, err error) {
				assert.Nil(t, category, "unexpected non-nil category")
				assert.Error(t, err, "unexpected nil-error")
			},
		},
		{
			name:          "fail-invalid-category-type-id",
			categoryName:  strings.Repeat("a", 255),
			categoryRefId: 199,
			assertions: func(t *testing.T, db *sqlx.DB, category *model.Category, err error) {
				assert.Nil(t, category, "unexpected non-nil category")
				assert.Error(t, err, "unexpected nil-error")
			},
		},
	}
}

func Test_AddCategory(t *testing.T) {
	for _, testCase := range getAddCategoryTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
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

			cat := &model.CreateCategory{
				Name:              testCase.categoryName,
				CategoryTypeRefId: testCase.categoryRefId,
			}

			createdCat, err := m.AddCategory(testCtx, txHandlerDb, cat)
			testCase.assertions(t, db, createdCat, err)
		})
	}
}

type testCaseDeleteCategoryTestCase struct {
	name       string
	id         int
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB)
}

func getDeleteCategoryTestCases() []testCaseDeleteCategoryTestCase {
	return []testCaseDeleteCategoryTestCase{
		{
			name: "success",
			id:   1,
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindCategory(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the category")

				assert.Equal(t, 0, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB) {
				entry := mysqlmodel.Category{
					CategoryTypeRefID: 1,
					Name:              "test",
				}
				err := entry.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")
			},
		},
	}
}

func Test_DeleteCategory(t *testing.T) {
	for _, testCase := range getDeleteCategoryTestCases() {
		db, cp, cleanup := mysqlhelper.TestGetMockMariaDB(t)
		t.Run(testCase.name, func(t *testing.T) {
			require.NotNil(t, testCase.mutations, "unexpected nil mutations")
			require.NotNil(t, testCase.assertions, "unexpected nil assertions")

			defer cleanup()
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

			testCase.mutations(t, db)
			err = m.DeleteCategory(testCtx, txHandlerDb, testCase.id)
			testCase.assertions(t, db, testCase.id, err)
		})
	}
}
