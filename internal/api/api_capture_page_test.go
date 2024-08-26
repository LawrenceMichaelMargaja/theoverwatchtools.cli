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
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type argsCreateCapturePage struct {
	User              mysqlmodel.User
	CapturePage       mysqlmodel.CapturePage
	CapturePageSet    mysqlmodel.CapturePageSet
	Organization      mysqlmodel.Organization
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
				User: mysqlmodel.User{
					ID:                4,
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					Firstname:         "Demby",
					Lastname:          "Abella",
					CategoryTypeRefID: 1,
				},
				Organization: mysqlmodel.Organization{
					ID:            1,
					Name:          "test",
					CreatedBy:     null.IntFrom(4),
					LastUpdatedBy: null.IntFrom(4),
				},
				CapturePageSet: mysqlmodel.CapturePageSet{
					ID:                1,
					Name:              "lawrence",
					CreatedBy:         null.IntFrom(1),
					LastUpdatedBy:     null.IntFrom(1),
					OrganizationRefID: null.IntFrom(1),
				},
				CapturePage: mysqlmodel.CapturePage{
					ID:            1,
					Name:          "Mohamed",
					CreatedBy:     null.IntFrom(1),
					LastUpdatedBy: null.IntFrom(1),
				},
				CreateCapturePage: &model.CreateCapturePage{
					Name:             "younes",
					UserId:           null.IntFrom(4).Int,
					CapturePageSetId: 1,
				},
			},
			mutations: func(t *testing.T, db *sqlx.DB, modules *testassets.Container, args *argsCreateCapturePage) {
				err := args.User.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the user db")

				err = args.Organization.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the organization db")

				err = args.CapturePageSet.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the CapturePageSet db")

				err = args.CapturePage.Insert(context.Background(), db, boil.Infer())
				require.NoError(t, err, "error inserting in the CapturePage db")
			},
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
