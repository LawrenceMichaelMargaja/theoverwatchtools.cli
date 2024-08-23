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

type argsGetOrganizations struct {
	User         mysqlmodel.User
	Organization []mysqlmodel.Organization
}

type testCaseGetOrganizations struct {
	name       string
	filter     *model.OrganizationFilters
	args       *argsGetOrganizations
	mutations  func(t *testing.T, db *sqlx.DB, args *argsGetOrganizations) []int
	assertions func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int)
}

func getTestCasesGetOrganizations() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name: "success-filter-ids-in",
			filter: &model.OrganizationFilters{
				IdsIn: []int{4},
			},
			args: &argsGetOrganizations{
				User: mysqlmodel.User{
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: []mysqlmodel.Organization{
					{
						ID:            4,
						Name:          "TEST",
						CreatedBy:     null.IntFrom(4),
						LastUpdatedBy: null.IntFrom(4),
						CreatedAt:     time.Now(),
						LastUpdatedAt: null.TimeFrom(time.Now()),
						IsActive:      true,
					},
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.Equal(t, len(ids), len(paginated.Organizations), "unexpected number of organizations returned")

				for _, org := range paginated.Organizations {
					assert.Contains(t, ids, org.Id, "returned organization ID not found in expected IDs")
				}

				assert.Equal(t, len(ids), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetOrganizations) []int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				var ids []int
				for _, organization := range args.Organization {
					err = organization.Insert(context.Background(), db, boil.Infer())
					require.NoError(t, err, "error inserting sample data")
					ids = append(ids, organization.ID)
				}

				return ids
			},
		},
		{
			name: "success-filter-names-in",
			filter: &model.OrganizationFilters{
				OrganizationNameIn: []string{"Lawrence", "Younes"},
			},
			args: &argsGetOrganizations{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: []mysqlmodel.Organization{
					{
						ID:            1,
						Name:          "Lawrence",
						CreatedBy:     null.IntFrom(4),
						LastUpdatedBy: null.IntFrom(4),
						CreatedAt:     time.Now(),
						LastUpdatedAt: null.TimeFrom(time.Now()),
						IsActive:      true,
					},
					{
						ID:            2,
						Name:          "Younes",
						CreatedBy:     null.IntFrom(4),
						LastUpdatedBy: null.IntFrom(4),
						CreatedAt:     time.Now(),
						LastUpdatedAt: null.TimeFrom(time.Now()),
						IsActive:      true,
					},
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.Equal(t, len(ids), len(paginated.Organizations), "unexpected number of organizations returned")

				for _, org := range paginated.Organizations {
					assert.Contains(t, ids, org.Id, "returned organization ID not found in expected IDs")
				}

				assert.Equal(t, len(ids), paginated.Pagination.RowCount, "unexpected row count")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetOrganizations) []int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting the user into the database")

				var ids []int
				for _, organization := range args.Organization {
					err = organization.Insert(context.Background(), db, boil.Infer())
					require.NoError(t, err, "error inserting sample data")
					ids = append(ids, organization.ID)
				}

				return ids
			},
		},
		{
			name: "success-multiple-filters",
			filter: &model.OrganizationFilters{
				OrganizationNameIn: []string{"Demby"},
				IdsIn:              []int{1},
				CreatedBy:          null.IntFrom(4),
			},
			args: &argsGetOrganizations{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: []mysqlmodel.Organization{
					{
						ID:            1,
						Name:          "Demby",
						CreatedBy:     null.IntFrom(4),
						LastUpdatedBy: null.IntFrom(4),
						CreatedAt:     time.Now(),
						LastUpdatedAt: null.TimeFrom(time.Now()),
						IsActive:      true,
					},
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetOrganizations) []int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				var ids []int
				for _, organization := range args.Organization {
					err = organization.Insert(context.Background(), db, boil.Infer())
					require.NoError(t, err, "error inserting sample data")
					ids = append(ids, organization.ID)
				}

				return ids
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil organizations")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 1, "unexpected greater than 1 organization")
				assert.True(t, paginated.Pagination.RowCount == 1, "unexpected count to be greater than 1 organization")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
		{
			name: "empty-results",
			filter: &model.OrganizationFilters{
				OrganizationNameIn:   []string{"Saul Goodman"},
				IdsIn:                []int{1},
				OrganizationIsActive: null.BoolFrom(false),
			},
			args: &argsGetOrganizations{},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsGetOrganizations) []int {
				return nil
			},
			assertions: func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int) {
				require.NoError(t, err, "unexpected error")
				require.NotNil(t, paginated, "unexpected nil paginated")
				require.NotNil(t, paginated.Organizations, "unexpected nil categories")
				require.NotNil(t, paginated.Pagination, "unexpected nil pagination")
				assert.True(t, len(paginated.Organizations) == 0, "unexpected greater than 1 category")
				assert.True(t, paginated.Pagination.RowCount == 0, "unexpected count to be greater than 1 category")

				modelhelpers.AssertNonEmptyOrganizations(t, paginated.Organizations)
			},
		},
	}
}

func Test_GetOrganizations(t *testing.T) {
	for _, testCase := range getTestCasesGetOrganizations() {
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
			paginated, err := m.GetOrganizations(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err, id)
		})
	}
}

type argsDeleteOrganization struct {
	User         mysqlmodel.User
	Organization mysqlmodel.Organization
}

type testCaseDeleteOrganization struct {
	name       string
	args       *argsDeleteOrganization
	assertions func(t *testing.T, db *sqlx.DB, id int)
	mutations  func(t *testing.T, db *sqlx.DB, args *argsDeleteOrganization)
}

func getDeleteOrganizationTestCases() []testCaseDeleteOrganization {
	return []testCaseDeleteOrganization{
		{
			name: "success",
			args: &argsDeleteOrganization{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "Lawrence",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Nil(t, err, "unexpected non-nil error")
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, false, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")
			},
		},
		{
			name: "failure-non-existent-organization",
			args: &argsDeleteOrganization{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "Lawrence",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int) {
				_, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.Error(t, err, "expected an error for non-existent organization")
				require.Contains(t, err.Error(), "no rows", "error should indicate that the organization was not found")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsDeleteOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")
			},
		},
	}
}

func Test_DeleteOrganization(t *testing.T) {
	for _, testCase := range getDeleteOrganizationTestCases() {
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

			testCase.mutations(t, db, testCase.args)
			err = m.DeleteOrganization(testCtx, txHandlerDb, testCase.args.Organization.ID)
			require.NoError(t, err, "error deleting organization")
			testCase.assertions(t, db, testCase.args.Organization.ID)
		})
	}
}

type argsUpdateOrganization struct {
	User                   mysqlmodel.User
	Organization           mysqlmodel.Organization
	UpdateOrganizationData *model.UpdateOrganization
}

type testCaseUpdateOrganization struct {
	name       string
	args       *argsUpdateOrganization
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB, args *argsUpdateOrganization)
}

func getUpdateOrganizationTestCases() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
		{
			name: "success",
			args: &argsUpdateOrganization{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				UpdateOrganizationData: &model.UpdateOrganization{
					Id:   4,
					Name: null.StringFrom("Organization A new"),
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")

				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
		},
		{
			name: "fail",
			args: &argsUpdateOrganization{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				UpdateOrganizationData: &model.UpdateOrganization{
					Id:   32123,
					Name: null.StringFrom("Wrong Name"),
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Error(t, err, "expected error when updating a non-existent organization")
				assert.Contains(t, err.Error(), "expected exactly one entry for entity: 32123")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsUpdateOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
		},
	}
}

func Test_UpdateOrganization(t *testing.T) {
	for _, testCase := range getUpdateOrganizationTestCases() {
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

		testCase.mutations(t, db, testCase.args)

		_, err = m.UpdateOrganization(testCtx, txHandlerDb, testCase.args.UpdateOrganizationData)
		testCase.assertions(t, db, testCase.args.UpdateOrganizationData.Id, err)
	}
}

type argsAddOrganization struct {
	User                   mysqlmodel.User
	Organization           mysqlmodel.Organization
	CreateOrganizationData *model.CreateOrganization
}

type testCaseCreateOrganization struct {
	name       string
	args       *argsAddOrganization
	assertions func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error)
	mutations  func(t *testing.T, db *sqlx.DB, args *argsAddOrganization)
}

func getAddOrganizationTestCases() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
		{
			name: "success",
			args: &argsAddOrganization{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CreateOrganizationData: &model.CreateOrganization{
					Name:   "Demby",
					UserId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsAddOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.NotNil(t, organization, "unexpected nil organization")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
		{
			name: "fail-name-too-long",
			args: &argsAddOrganization{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					Email:             "demby@test.com",
					Password:          "password",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            4,
					Name:          "TEST",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
					CreatedAt:     time.Now(),
					LastUpdatedAt: null.TimeFrom(time.Now()),
					IsActive:      true,
				},
				CreateOrganizationData: &model.CreateOrganization{
					Name:   strings.Repeat("DembyYounesLawrence", 256),
					UserId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsAddOrganization) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")
			},
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.Nil(t, organization, "unexpected non-nil organization")
				assert.Error(t, err, "expected an error due to name exceeding length limit")
			},
		},
	}
}

func Test_AddOrganization(t *testing.T) {
	for _, testCase := range getAddOrganizationTestCases() {
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

			createOrganization, err := m.AddOrganization(testCtx, txHandlerDb, testCase.args.CreateOrganizationData)
			testCase.assertions(t, db, createOrganization, err)
		})
	}
}

type argsRestoreOrganization struct {
	Organization mysqlmodel.Organization
}

type testCaseRestoreOrganization struct {
	name       string
	args       *argsRestoreOrganization
	assertions func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations  func(t *testing.T, db *sqlx.DB, args *argsRestoreOrganization)
}

func getRestoreOrganizationTestCases() []testCaseRestoreOrganization {
	return []testCaseRestoreOrganization{
		{
			name: "success",
			args: &argsRestoreOrganization{
				Organization: mysqlmodel.Organization{
					ID:       1,
					Name:     "Organization A",
					IsActive: false,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				require.Nil(t, err, "unexpected non-nil error")
				entry, err := mysqlmodel.FindOrganization(context.TODO(), db, id)
				require.NoError(t, err, "unexpected error fetching the organization")
				assert.Equal(t, true, entry.IsActive)
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreOrganization) {
				err := args.Organization.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "unexpected insert error")
			},
		},
		{
			name: "fail-missing-entry-to-update",
			args: &argsRestoreOrganization{
				Organization: mysqlmodel.Organization{
					ID:       1,
					Name:     "Non-existent Organization",
					IsActive: false,
				},
			},
			assertions: func(t *testing.T, db *sqlx.DB, id int, err error) {
				_, fetchErr := mysqlmodel.FindOrganization(context.TODO(), db, id)
				assert.Error(t, fetchErr, "organization should not exist")
			},
			mutations: func(t *testing.T, db *sqlx.DB, args *argsRestoreOrganization) {},
		},
	}
}

func Test_RestoreOrganization(t *testing.T) {
	for _, testCase := range getRestoreOrganizationTestCases() {
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

			testCase.mutations(t, db, testCase.args)
			err = m.RestoreOrganization(testCtx, txHandlerDb, testCase.args.Organization.ID)
			testCase.assertions(t, db, testCase.args.Organization.ID, err)
		})
	}
}
