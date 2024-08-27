package router

import (
	"context"
	"errors"
	"gomap/src/router"
	"gomap/src/testUtils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateMapUIHandlerIntegration(t *testing.T) {
	timeout := 5 * time.Second

	testCases := []struct {
		testName          string
		sheetId           string
		mockCsvResponse   string
		mockCsvStatusCode int
		expectedStatus    int
		expectedMessage   string
		mockRedisSetup    func(mockRedisClient *testUtils.MockRedisClient, ctx context.Context)
	}{
		{
			testName:          "no sheetId provided",
			sheetId:           "",
			mockCsvResponse:   "",
			mockCsvStatusCode: http.StatusOK,
			expectedStatus:    http.StatusBadRequest,
			expectedMessage:   "Missing sheetId parameter",
			mockRedisSetup:    nil,
		},
		{
			testName:          "fetch and parse locations successfully",
			sheetId:           "mockSheetId",
			mockCsvResponse:   "",
			mockCsvStatusCode: http.StatusOK,
			expectedStatus:    http.StatusOK,
			expectedMessage:   `<iframe src="/?sheetId=mockSheetId" width="100%" height="300" frameborder="0"></iframe>`,
			mockRedisSetup: func(mockRedisClient *testUtils.MockRedisClient, ctx context.Context) {
				mockRedisClient.On("Set",
					mock.Anything,
					"mockSheetId",
					mock.Anything,
					mock.AnythingOfType("time.Duration"),
				).Return(&redis.StatusCmd{})
			},
		},
		{
			testName:          "fetch locations fails",
			sheetId:           "mockSheetId",
			mockCsvResponse:   "",
			mockCsvStatusCode: http.StatusInternalServerError,
			expectedStatus:    http.StatusInternalServerError,
			expectedMessage:   "failed to fetch CSV data",
			mockRedisSetup:    nil,
		},
		{
			testName: "parse locations fails",
			sheetId:  "mockSheetId",
			mockCsvResponse: `Name,Address,City,State,Country,Website,Phone Number,Latitude,Longitude
        Location 1,Address 1,City 1,State 1,Country 1,https://example.com,1234567890,INVALID_LAT,INVALID_LONG`,
			mockCsvStatusCode: http.StatusOK,
			expectedStatus:    http.StatusInternalServerError,
			expectedMessage:   "failed to parse locations",
			mockRedisSetup:    nil,
		},
		{
			testName:          "process locations fails",
			sheetId:           "mockSheetId",
			mockCsvResponse:   "",
			mockCsvStatusCode: http.StatusOK,
			expectedStatus:    http.StatusInternalServerError,
			expectedMessage:   "failed to cache locations",
			mockRedisSetup: func(mockRedisClient *testUtils.MockRedisClient, ctx context.Context) {
				statusCmd := redis.NewStatusCmd(ctx)
				statusCmd.SetErr(errors.New("failed to cache locations"))
				mockRedisClient.On("Set",
					mock.Anything,
					"mockSheetId",
					mock.Anything,
					mock.Anything,
				).Return(statusCmd)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			mockRedisClient := &testUtils.MockRedisClient{}
			mockCSVServer := createMockCSVServer(testCase.mockCsvResponse, testCase.mockCsvStatusCode)
			defer mockCSVServer.Close()

			if testCase.mockRedisSetup != nil {
				testCase.mockRedisSetup(mockRedisClient, ctx)
			}

			r := router.InitRouter(router.RouterConfig{
				RedisClient:        mockRedisClient,
				Ctx:                ctx,
				BaseSpreadsheetUrl: mockCSVServer.URL + "?sheetId=%s",
			})

			form := strings.NewReader("sheetId=" + testCase.sheetId)
			req, err := http.NewRequest("POST", "/updateMapUI", form)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			routerRecorder := httptest.NewRecorder()
			r.ServeHTTP(routerRecorder, req)

			assert.Equal(t, testCase.expectedStatus, routerRecorder.Code)
			assert.Contains(t, routerRecorder.Body.String(), testCase.expectedMessage)

			defer func() {
				mockRedisClient.AssertExpectations(t)
				mockRedisClient.Calls = nil
				mockRedisClient.ExpectedCalls = nil
			}()
		})
	}
}

func createMockCSVServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))
}
