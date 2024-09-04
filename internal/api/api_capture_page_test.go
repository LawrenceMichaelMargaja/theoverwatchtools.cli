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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type argsCreateCapturePage struct {
	User              mysqlmodel.User
	CapturePage       mysqlmodel.CapturePage
	CapturePageSet    mysqlmodel.CapturePageSet
	CreateCapturePage *model.CreateCapturePage
}

type testCaseCreateCapturePage struct {
	name              string
	fnGetTestServices func(t *testing.T) (*testassets.Container, func())
	args              *argsCreateCapturePage
	mutations         func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsCreateCapturePage)
	assertions        func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesCreateCapturePage() []testCaseCreateCapturePage {
	return []testCaseCreateCapturePage{
		{
			name: "success",
			fnGetTestServices: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			args: &argsCreateCapturePage{
				CreateCapturePage: &model.CreateCapturePage{
					Name:             "younes",
					UserId:           null.IntFrom(2).Int,
					CapturePageSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsCreateCapturePage) {},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				respStr := string(resp)
				require.NotNilf(t, resp, "unexpected nil response: %s", respStr)
				require.Equal(t, http.StatusCreated, respCode, "unexpected non-equal response code: %s", respStr)

				var capturePage *model.CapturePage
				err := json.Unmarshal(resp, &capturePage)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				require.NotNil(t, capturePage, "unexpected nil capture page")

				modelhelpers.AssertNonEmptyCapturePages(t, []model.CapturePage{*capturePage})
			},
		},
	}
}

func Test_CreateCapturePage(t *testing.T) {
	for _, testCase := range getTestCasesCreateCapturePage() {
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
				Logger:              logger.New(context.TODO()),
			}

			testCase.mutations(t, db, handlers, testCase.args)

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			reqB, err := json.Marshal(testCase.args.CreateCapturePage)
			require.NoError(t, err, "unexpected error marshalling parameters")

			req := httptest.NewRequest(http.MethodPost, "/api/v1/capturepage", bytes.NewBuffer(reqB))

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

type argsListCapturePages struct {
	User              mysqlmodel.User
	CreateCapturePage *model.CreateCapturePage
	Category          mysqlmodel.Category
	CapturePage       mysqlmodel.CapturePage
	CapturePageSet    mysqlmodel.CapturePageSet
	Organization      mysqlmodel.Organization
}

type testCaseListCapturePages struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	args            *argsListCapturePages
	mutations       func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsListCapturePages)
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesListCapturePages() []testCaseListCapturePages {
	testCases := []testCaseListCapturePages{
		{
			name: "success",
			queryParameters: map[string]interface{}{
				"ids_in": []int{1},
			},
			args: &argsListCapturePages{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					CategoryTypeRefID: 2,
				},
				CapturePage: mysqlmodel.CapturePage{
					ID:               2,
					Name:             "demby",
					CreatedBy:        null.IntFrom(4),
					LastUpdatedBy:    null.IntFrom(4),
					CapturePageSetID: 3,
				},
				Organization: mysqlmodel.Organization{
					ID:            5,
					Name:          "demby",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				CreateCapturePage: &model.CreateCapturePage{
					Name:             "demby",
					UserId:           1,
					CapturePageSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsListCapturePages) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the organization db")

				err = args.CapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the CapturePage db")

				_, err = modules.CapturePageService.AddCapturePage(context.Background(), args.CreateCapturePage)
				require.NoError(t, err, "error adding the capture page")
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				assert.Equal(t, http.StatusOK, respCode, "unexpected non-equal response code")

				var respPaginated model.PaginatedCapturePages
				err := json.Unmarshal(resp, &respPaginated)
				require.NoError(t, err, "unexpected error unmarshalling the response")

				assert.Len(t, respPaginated.CapturePages, 1, "unexpected number of capture pages")
				assert.Equal(t, 1, respPaginated.Pagination.RowCount, "unexpected row_count")
				assert.Equal(t, 1, respPaginated.Pagination.TotalCount, "unexpected total_count")
				assert.Equal(t, 1, respPaginated.Pagination.Page, "unexpected page")
			},
		},
	}

	return testCases
}

func Test_ListCapturePages(t *testing.T) {
	for _, testCase := range getTestCasesListCapturePages() {
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
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			url := strutil.AppendQueryToURL("/api/v1/capturepage", testCase.queryParameters)
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

type argsDeleteCapturePages struct {
	User              mysqlmodel.User
	Category          mysqlmodel.Category
	CapturePage       mysqlmodel.CapturePage
	CapturePageSet    mysqlmodel.CapturePageSet
	Organization      mysqlmodel.Organization
	CreateCapturePage *model.CreateCapturePage
}

type testCaseDeleteCapturePages struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	args            *argsDeleteCapturePages
	mutations       func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsDeleteCapturePages) int
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesDeleteCapturePages() []testCaseDeleteCapturePages {
	testCases := []testCaseDeleteCapturePages{
		{
			name: "success",
			queryParameters: map[string]interface{}{
				"ids_in": []int{1},
			},
			args: &argsDeleteCapturePages{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					CategoryTypeRefID: 2,
				},
				CapturePage: mysqlmodel.CapturePage{
					ID:               2,
					Name:             "demby",
					CreatedBy:        null.IntFrom(4),
					LastUpdatedBy:    null.IntFrom(4),
					CapturePageSetID: 3,
				},
				Organization: mysqlmodel.Organization{
					ID:            5,
					Name:          "demby",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				CreateCapturePage: &model.CreateCapturePage{
					Name:             "demby",
					UserId:           1,
					CapturePageSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsDeleteCapturePages) int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the organization db")

				err = args.CapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the CapturePage db")

				createdCapturePage, addCapturePageError := modules.CapturePageService.AddCapturePage(context.Background(), args.CreateCapturePage)
				require.NoError(t, addCapturePageError, "error adding the capture page")

				return createdCapturePage.Id
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusNoContent, respCode, "unexpected response code")
				assert.Empty(t, resp, "expected empty response body for no-content response")
			},
		},
	}

	return testCases
}

func Test_DeleteCapturePage(t *testing.T) {
	for _, testCase := range getTestCasesDeleteCapturePages() {
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
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			id := testCase.mutations(t, db, handlers, testCase.args)

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/capturepage/"+strconv.Itoa(id), nil)
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

type argsRestoreCapturePages struct {
	User              mysqlmodel.User
	Category          mysqlmodel.Category
	CapturePage       mysqlmodel.CapturePage
	CapturePageSet    mysqlmodel.CapturePageSet
	Organization      mysqlmodel.Organization
	CreateCapturePage *model.CreateCapturePage
}

type testCaseRestoreCapturePages struct {
	name            string
	getContainer    func(t *testing.T) (*testassets.Container, func())
	args            *argsRestoreCapturePages
	mutations       func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsRestoreCapturePages) int
	queryParameters map[string]interface{}
	assertions      func(t *testing.T, resp []byte, respCode int, capID int, db *sqlx.DB)
}

func getTestCasesRestoreCapturePages() []testCaseRestoreCapturePages {
	testCases := []testCaseRestoreCapturePages{
		{
			name: "success",
			queryParameters: map[string]interface{}{
				"ids_in": []int{1},
			},
			args: &argsRestoreCapturePages{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					CreatedBy:         null.IntFrom(4),
					LastUpdatedBy:     null.IntFrom(4),
					CategoryTypeRefID: 2,
				},
				CapturePage: mysqlmodel.CapturePage{
					ID:               2,
					Name:             "demby",
					CreatedBy:        null.IntFrom(4),
					LastUpdatedBy:    null.IntFrom(4),
					CapturePageSetID: 3,
					IsActive:         false,
				},
				Organization: mysqlmodel.Organization{
					ID:            5,
					Name:          "demby",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsRestoreCapturePages) int {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the organization db")

				err = args.CapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the CapturePage db")

				return args.CapturePage.ID
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int, capID int, db *sqlx.DB) {
				require.Equal(t, http.StatusNoContent, respCode, "unexpected response code")
				assert.Empty(t, resp, "unexpected non-empty response body")

				capturePage, err := mysqlmodel.FindCapturePage(context.TODO(), db, capID)
				require.NoError(t, err, "unexpected error fetching organization from database")

				assert.True(t, capturePage.IsActive, "expected is_active to be true")
			},
		},
	}

	return testCases
}

func Test_RestoreCapturePage(t *testing.T) {
	for _, testCase := range getTestCasesRestoreCapturePages() {
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
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			id := testCase.mutations(t, db, handlers, testCase.args)

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/capturepage/"+strconv.Itoa(id), nil)
			req.Header = map[string][]string{
				"Content-Type":    {"application/json"},
				"Accept-Encoding": {"gzip", "deflate", "br"},
			}

			resp, err := api.app.Test(req, 100)
			require.NoError(t, err, "unexpected error executing test")

			respBytes, err := io.ReadAll(resp.Body)
			require.Nil(t, err, "unexpected error reading the response")
			testCase.assertions(t, respBytes, resp.StatusCode, id, db)
		})
	}
}

type argsUpdateCapturePages struct {
	User                  mysqlmodel.User
	Category              mysqlmodel.Category
	CapturePage           mysqlmodel.CapturePage
	CapturePageSet        mysqlmodel.CapturePageSet
	Organization          mysqlmodel.Organization
	UpdateCapturePageData *model.UpdateCapturePage
	CreateCapturePageData *model.CreateCapturePage
}

type testCaseUpdateCapturePages struct {
	name         string
	getContainer func(t *testing.T) (*testassets.Container, func())
	args         *argsUpdateCapturePages
	mutations    func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsUpdateCapturePages)
	assertions   func(t *testing.T, resp []byte, respCode int)
}

func getTestCasesUpdateCapturePages() []testCaseUpdateCapturePages {
	testCases := []testCaseUpdateCapturePages{
		{
			name: "success",
			args: &argsUpdateCapturePages{
				User: mysqlmodel.User{
					ID:                4,
					Firstname:         "Demby",
					Lastname:          "Abella",
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					CategoryTypeRefID: 1,
				},
				CapturePage: mysqlmodel.CapturePage{
					ID:               1,
					Name:             "demby",
					CreatedBy:        null.IntFrom(4),
					LastUpdatedBy:    null.IntFrom(4),
					CapturePageSetID: 1,
					IsActive:         false,
				},
				UpdateCapturePageData: &model.UpdateCapturePage{
					Id: 1,
					Name: null.String{
						String: "lawrence",
						Valid:  true,
					},
					UserId:           null.IntFrom(2),
					CapturePageSetId: 1,
				},
				CreateCapturePageData: &model.CreateCapturePage{
					Name:             "demby",
					UserId:           2,
					CapturePageSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsUpdateCapturePages) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				_, err = modules.CapturePageService.AddCapturePage(context.Background(), args.CreateCapturePageData)
				require.NoError(t, err, "error adding the capture page")
			},
			getContainer: func(t *testing.T) (*testassets.Container, func()) {
				ctn, cleanup := testassets.GetConcreteContainer(t)
				return ctn, func() {
					cleanup()
				}
			},
			assertions: func(t *testing.T, resp []byte, respCode int) {
				require.Equal(t, http.StatusOK, respCode)
				var capturePage *model.CapturePage
				err := json.Unmarshal(resp, &capturePage)
				require.NoError(t, err, "unexpected error unmarshalling the response")
				modelhelpers.AssertNonEmptyCapturePages(t, []model.CapturePage{*capturePage})
			},
		},
	}

	return testCases
}

func Test_UpdateCapturePage(t *testing.T) {
	for _, testCase := range getTestCasesUpdateCapturePages() {
		t.Run(testCase.name, func(t *testing.T) {
			db, _, cleanup := mysqlhelper.TestGetMockMariaDB(t)
			defer cleanup()

			handlers, _ := testCase.getContainer(t)

			cfg := &Config{
				BaseUrl:             testassets.MockBaseUrl,
				Port:                3000,
				CategoryService:     handlers.CategoryService,
				OrganizationService: handlers.OrganizationService,
				CapturePageService:  handlers.CapturePageService,
				Logger:              logger.New(context.TODO()),
			}

			api, err := New(cfg)
			require.NoError(t, err, "unexpected error instantiating api")
			require.NotNil(t, api, "unexpected api nil instance")

			testCase.mutations(t, db, handlers, testCase.args)

			reqB, err := json.Marshal(testCase.args.UpdateCapturePageData)
			require.NoError(t, err, "unexpected error marshalling parameters")

			req := httptest.NewRequest(http.MethodPatch, "/api/v1/capturepage", bytes.NewBuffer(reqB))
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
