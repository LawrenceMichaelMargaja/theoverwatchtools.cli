package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dembygenesis/local.tools/internal/api/testassets"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/model/modelhelpers"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/assets/mysqlmodel"
	"github.com/dembygenesis/local.tools/internal/persistence/database_helpers/mysql/mysqlhelper"
	"github.com/dembygenesis/local.tools/internal/utilities/strutil"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type argsCreateClickTracker struct {
	User               mysqlmodel.User
	ClickTracker       mysqlmodel.ClickTracker
	ClickTrackerSet    mysqlmodel.ClickTrackerSet
	CreateClickTracker *model.CreateClickTracker
}

type testCaseCreateClickTracker struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testassets.Container, func())
	args              *argsCreateClickTracker
	mutations         func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsCreateClickTracker)
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesCreateClickTracker() []testCaseCreateClickTracker {
	return []testCaseCreateClickTracker{
		{
			name: "success",
			fnGetTestServices: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			args: &argsCreateClickTracker{
				CreateClickTracker: &model.CreateClickTracker{
					Name:              "younes",
					UserId:            null.IntFrom(1).Int,
					ClickTrackerSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsCreateClickTracker) {
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				respStr := string(resp)
				require.NotNilf(t, resp, "unexpected nil response: %s", respStr)
				require.Equal(t, http.StatusCreated, respCode, "unexpected non-equal response code: %s", respStr)

				var ClickTracker *model.ClickTracker
				err := json.Unmarshal(resp, &ClickTracker)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				require.NotNil(t, ClickTracker, "unexpected nil click tracker")

				modelhelpers.AssertNonEmptyClickTrackers(t, []model.ClickTracker{*ClickTracker})
			},
		},
		{
			name: "fail-mock-server-error",
			fnGetTestServices: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			args: &argsCreateClickTracker{
				CreateClickTracker: &model.CreateClickTracker{
					Name: "younes",
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsCreateClickTracker) {
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.NotNil(t, resp, "unexpected nil response")
				assert.Equal(t, http.StatusInternalServerError, respCode)
			},
		},
	}
}

func Test_CreateClickTracker(t *testing.T) {
	for _, testCase := range getTestCasesCreateClickTracker() {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()

			handlers, _ := testCase.fnGetTestServices(t)

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				OrganizationService: handlers.OrganizationService,
				CapturePageService:  handlers.CapturePageService,
				ClickTrackerService: handlers.ClickTrackerService,
				Logger:              logger.New(context.TODO()),
			}

			testCase.mutations(t, db, handlers, testCase.args)

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(testCase.args.CreateClickTracker)
			require.NoError(t, err, "unexpected error marshalling parameters")

			req := httptest.NewRequest(http.MethodPost, "/api/v1/clicktracker", bytes.NewBuffer(reqB))

			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode)
		})
	}
}

type argsListClickTrackers struct {
	CreateClickTracker      *model.CreateClickTracker
	CreateClickTrackerTwo   *model.CreateClickTracker
	CreateClickTrackerThree *model.CreateClickTracker
}

type testCaseClickTrackers struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	args            *argsListClickTrackers
	mutations       func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsListClickTrackers)
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesClickTrackers() []testCaseClickTrackers {
	testCases := []testCaseClickTrackers{
		{
			name: "success",
			queryParameters: map[string]interface{}{
				"ids_in": []int{1},
			},
			args: &argsListClickTrackers{
				CreateClickTracker: &model.CreateClickTracker{
					Name:              "Younes",
					UserId:            1,
					ClickTrackerSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsListClickTrackers) {
				_, err := modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTracker)
				require.NoError(t, err, "error adding the click tracker")

			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")

				var respPaginated model.PaginatedClickTrackers
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				assert.Len(t, respPaginated.ClickTrackers, 1, "unexpected number of click trackers")
				assert.Equal(t, 1, respPaginated.Pagination.RowCount, "unexpected row_count")
				assert.Equal(t, 1, respPaginated.Pagination.TotalCount, "unexpected total_count")
				assert.Equal(t, 1, respPaginated.Pagination.Page, "unexpected page")
			},
		},
		{
			name: "success-limit-offset-page-1",
			queryParameters: map[string]interface{}{
				"ids_in":   []int{1, 2, 3},
				"page":     1,
				"max_rows": 1,
			},
			args: &argsListClickTrackers{
				CreateClickTracker: &model.CreateClickTracker{
					Name:              "Younes",
					UserId:            1,
					ClickTrackerSetId: 1,
				},
				CreateClickTrackerTwo: &model.CreateClickTracker{
					Name:              "demby",
					UserId:            2,
					ClickTrackerSetId: 1,
				},
				CreateClickTrackerThree: &model.CreateClickTracker{
					Name:              "lawrence",
					UserId:            3,
					ClickTrackerSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsListClickTrackers) {
				_, err := modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTracker)
				require.NoError(t, err, "error adding the click tracker")

				_, err = modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTrackerTwo)
				require.NoError(t, err, "error adding the click tracker")

				_, err = modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTrackerThree)
				require.NoError(t, err, "error adding the click tracker")
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equalf(t, http.StatusOK, respCode, "unexpected status code: %v, with resp: %s", respCode, string(resp))

				var respPaginated model.PaginatedClickTrackers
				err := json.Unmarshal(resp, &respPaginated)

				require.NoError(t, err, "unexpected error unmarshalling the response")
				assert.NotNil(t, respPaginated.ClickTrackers, "unexpected empty click trackers")
				assert.Greater(t, respPaginated.Pagination.TotalCount, 1, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 1, "unexpected page")
			},
		},
		{
			name: "success-limit-offset-page-2",
			queryParameters: map[string]interface{}{
				"ids_in":   []int{1, 2, 3},
				"page":     2,
				"max_rows": 2,
			},
			args: &argsListClickTrackers{
				CreateClickTracker: &model.CreateClickTracker{
					Name:              "Younes",
					UserId:            1,
					ClickTrackerSetId: 1,
				},
				CreateClickTrackerTwo: &model.CreateClickTracker{
					Name:              "demby",
					UserId:            2,
					ClickTrackerSetId: 1,
				},
				CreateClickTrackerThree: &model.CreateClickTracker{
					Name:              "lawrence",
					UserId:            3,
					ClickTrackerSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsListClickTrackers) {
				_, err := modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTracker)
				require.NoError(t, err, "error adding the click tracker")

				_, err = modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTrackerTwo)
				require.NoError(t, err, "error adding the click tracker")

				_, err = modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTrackerThree)
				require.NoError(t, err, "error adding the click tracker")
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equalf(t, http.StatusOK, respCode, "unexpected status code: %v, with resp: %s", respCode, string(resp))

				var respPaginated model.PaginatedClickTrackers
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				assert.NotNil(t, respPaginated.ClickTrackers, "unexpected empty click trackers")
				assert.True(t, respPaginated.Pagination.TotalCount > 2, "unexpected total_count")
				assert.Equal(t, respPaginated.Pagination.Page, 2, "unexpected page")
			},
		},
		{
			name: "success-all-filters",
			queryParameters: map[string]interface{}{
				"ids_in":                  []int{1, 2, 3},
				"click_tracker_set_id_in": []int{1},
				"click_tracker_name_in":   []string{"Admin"},
			},
			args: &argsListClickTrackers{
				CreateClickTracker: &model.CreateClickTracker{
					Name:              "Admin",
					UserId:            1,
					ClickTrackerSetId: 1,
				},
				CreateClickTrackerTwo: &model.CreateClickTracker{
					Name:              "demby",
					UserId:            2,
					ClickTrackerSetId: 1,
				},
				CreateClickTrackerThree: &model.CreateClickTracker{
					Name:              "lawrence",
					UserId:            3,
					ClickTrackerSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsListClickTrackers) {
				_, err := modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTracker)
				require.NoError(t, err, "error adding the click tracker")

				_, err = modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTrackerTwo)
				require.NoError(t, err, "error adding the click tracker")

				_, err = modules.ClickTrackerService.AddClickTracker(context.Background(), args.CreateClickTrackerThree)
				require.NoError(t, err, "error adding the click tracker")
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				var respPaginated model.PaginatedClickTrackers
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")
				assert.True(t, len(respPaginated.ClickTrackers) > 0, "unexpected empty click trackers")
				assert.True(t, respPaginated.Pagination.MaxRows > 0, "unexpected empty rows")
				assert.True(t, respPaginated.Pagination.RowCount > 0, "unexpected empty count")
				assert.True(t, len(respPaginated.Pagination.Pages) > 0, "unexpected empty pages")
			},
		},
		{
			name:            "empty_store",
			queryParameters: map[string]interface{}{},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsListClickTrackers) {
				store := modules.MySQLStore
				require.NotNil(t, store, "unexpected nil: store")

				connProvider := modules.ConnProvider
				require.NotNil(t, store, "unexpected nil: txProvider")

				tx, err := connProvider.Tx(context.TODO())
				require.NoError(t, err, "unexpected err for getting tx")

				err = store.DropClickTrackerTable(context.TODO(), tx)
				require.NoError(t, err, "unexpected err for drop command")

				err = tx.Commit(context.TODO())
				require.NoError(t, err, "unexpected err on commit")
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				assert.Equal(t, http.StatusInternalServerError, respCode)
			},
		},
	}

	return testCases
}

func Test_ListClickTrackers(t *testing.T) {
	for _, testCase := range getTestCasesClickTrackers() {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()
			if testCase.queryParameters == nil {
				testCase.queryParameters = make(map[string]interface{})
			}

			handlers, _ := testCase.getContainer(t)

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				OrganizationService: handlers.OrganizationService,
				CapturePageService:  handlers.CapturePageService,
				ClickTrackerService: handlers.ClickTrackerService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := strutil.AppendQueryToURL("/api/v1/clicktracker", testCase.queryParameters)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			testCase.mutations(t, db, handlers, testCase.args)

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode)
		})
	}
}
