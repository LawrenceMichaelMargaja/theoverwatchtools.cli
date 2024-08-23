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

type argsMutationsGetOrganizations struct {
	User         mysqlmodel.User
	Organization []mysqlmodel.Organization
}

type testCaseGetOrganizations struct {
	name          string
	filter        *model.OrganizationFilters
	argsMutations *argsMutationsGetOrganizations
	mutations     func(t *testing.T, db *sqlx.DB, args *argsMutationsGetOrganizations) []int
	assertions    func(t *testing.T, db *sqlx.DB, paginated *model.PaginatedOrganizations, err error, ids []int)
}

func getTestCasesGetOrganizations() []testCaseGetOrganizations {
	return []testCaseGetOrganizations{
		{
			name: "success-filter-ids-in",
			filter: &model.OrganizationFilters{
				IdsIn: []int{4},
			},
			argsMutations: &argsMutationsGetOrganizations{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsGetOrganizations) []int {
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
			argsMutations: &argsMutationsGetOrganizations{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsGetOrganizations) []int {
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
			argsMutations: &argsMutationsGetOrganizations{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsGetOrganizations) []int {
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsGetOrganizations) []int {
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

			id := testCase.mutations(t, db, testCase.argsMutations)
			paginated, err := m.GetOrganizations(testCtx, txHandlerDb, testCase.filter)
			testCase.assertions(t, db, paginated, err, id)
		})
	}
}

type argsMutationsDeleteOrganization struct {
	User         mysqlmodel.User
	Organization mysqlmodel.Organization
}

type testCaseDeleteOrganization struct {
	name          string
	argsMutations *argsMutationsDeleteOrganization
	assertions    func(t *testing.T, db *sqlx.DB, id int)
	mutations     func(t *testing.T, db *sqlx.DB, args *argsMutationsDeleteOrganization) int
}

func getDeleteOrganizationTestCases() []testCaseDeleteOrganization {
	return []testCaseDeleteOrganization{
		{
			name: "success",
			argsMutations: &argsMutationsDeleteOrganization{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsDeleteOrganization) int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.TODO(), db, boil.Infer())
				assert.NoError(t, err, "unexpected insert error")

				return args.Organization.ID
			},
		},
		{
			name: "failure-non-existent-organization",
			argsMutations: &argsMutationsDeleteOrganization{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsDeleteOrganization) int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				return args.User.ID
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

			id := testCase.mutations(t, db, testCase.argsMutations)
			err = m.DeleteOrganization(testCtx, txHandlerDb, id)
			require.NoError(t, err, "error deleting organization")
			testCase.assertions(t, db, id)
		})
	}
}

type argsMutationsUpdateOrganization struct {
	User                   mysqlmodel.User
	Organization           mysqlmodel.Organization
	UpdateOrganizationData *model.UpdateOrganization
}

type testCaseUpdateOrganization struct {
	name          string
	argsMutations *argsMutationsUpdateOrganization
	assertions    func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations     func(t *testing.T, db *sqlx.DB, args *argsMutationsUpdateOrganization) *model.UpdateOrganization
}

func getUpdateOrganizationTestCases() []testCaseUpdateOrganization {
	return []testCaseUpdateOrganization{
		{
			name: "success",
			argsMutations: &argsMutationsUpdateOrganization{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsUpdateOrganization) *model.UpdateOrganization {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				return args.UpdateOrganizationData
			},
		},
		{
			name: "fail",
			argsMutations: &argsMutationsUpdateOrganization{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsUpdateOrganization) *model.UpdateOrganization {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				return args.UpdateOrganizationData
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

		updateData := testCase.mutations(t, db, testCase.argsMutations)

		_, err = m.UpdateOrganization(testCtx, txHandlerDb, updateData)
		testCase.assertions(t, db, updateData.Id, err)
	}
}

type argsMutationsAddOrganization struct {
	User                   mysqlmodel.User
	Organization           mysqlmodel.Organization
	CreateOrganizationData *model.CreateOrganization
}

type testCaseCreateOrganization struct {
	name          string
	argsMutations *argsMutationsAddOrganization
	assertions    func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error)
	mutations     func(t *testing.T, db *sqlx.DB, args *argsMutationsAddOrganization) *model.CreateOrganization
}

func getAddOrganizationTestCases() []testCaseCreateOrganization {
	return []testCaseCreateOrganization{
		{
			name: "success",
			argsMutations: &argsMutationsAddOrganization{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsAddOrganization) *model.CreateOrganization {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				return args.CreateOrganizationData
			},
			assertions: func(t *testing.T, db *sqlx.DB, organization *model.Organization, err error) {
				assert.NotNil(t, organization, "unexpected nil organization")
				assert.NoError(t, err, "unexpected non-nil error")

				modelhelpers.AssertNonEmptyOrganizations(t, []model.Organization{*organization})
			},
		},
		{
			name: "fail-name-too-long",
			argsMutations: &argsMutationsAddOrganization{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsAddOrganization) *model.CreateOrganization {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting sample data")

				return args.CreateOrganizationData
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

			organizationData := testCase.mutations(t, db, testCase.argsMutations)

			createOrganization, err := m.AddOrganization(testCtx, txHandlerDb, organizationData)
			testCase.assertions(t, db, createOrganization, err)
		})
	}
}

type argsMutationsRestoreOrganization struct {
	Organization mysqlmodel.Organization
}

type testCaseRestoreOrganization struct {
	name          string
	argsMutations *argsMutationsRestoreOrganization
	assertions    func(t *testing.T, db *sqlx.DB, id int, err error)
	mutations     func(t *testing.T, db *sqlx.DB, args *argsMutationsRestoreOrganization) (id int)
}

func getRestoreOrganizationTestCases() []testCaseRestoreOrganization {
	return []testCaseRestoreOrganization{
		{
			name: "success",
			argsMutations: &argsMutationsRestoreOrganization{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsRestoreOrganization) (id int) {
				err := args.Organization.Insert(context.TODO(), db, boil.Infer())
				require.NoError(t, err, "unexpected insert error")

				return args.Organization.ID
			},
		},
		{
			name: "fail-missing-entry-to-update",
			argsMutations: &argsMutationsRestoreOrganization{
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
			mutations: func(t *testing.T, db *sqlx.DB, args *argsMutationsRestoreOrganization) (id int) {
				return args.Organization.ID
			},
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

			id := testCase.mutations(t, db, testCase.argsMutations)
			err = m.RestoreOrganization(testCtx, txHandlerDb, id)
			testCase.assertions(t, db, id, err)
		})
	}
}
