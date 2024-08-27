package router

import (
	"context"
	"errors"
	"fmt"
	"gomap/src/router"
	"gomap/src/testUtils"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestCaseConfig struct {
	testName          string
	sheetId           string
	mockCsvResponse   string
	mockCsvStatusCode int
	expectedStatus    int
	expectedMessage   string
	mockRedisSetup    func(mockRedisClient *testUtils.MockRedisClient, ctx context.Context)
}

func TestUpdateMapUIHandlerIntegration(t *testing.T) {

	testCases := []TestCaseConfig{
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
			testName:          "no sheetId provided",
			sheetId:           "",
			mockCsvResponse:   "",
			mockCsvStatusCode: http.StatusOK,
			expectedStatus:    http.StatusBadRequest,
			expectedMessage:   "Missing sheetId parameter",
			mockRedisSetup:    nil,
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
			testName:          "invalid CSV format",
			sheetId:           "mockSheetId",
			mockCsvResponse:   "Name,Address\nLocation 1,Address 1\n",
			mockCsvStatusCode: http.StatusOK,
			expectedStatus:    http.StatusInternalServerError,
			expectedMessage:   "failed to parse locations",
			mockRedisSetup:    nil,
		},
		// {
		// 	testName:          "empty CSV",
		// 	sheetId:           "mockSheetId",
		// 	mockCsvResponse:   "",
		// 	mockCsvStatusCode: http.StatusOK,
		// 	expectedStatus:    http.StatusInternalServerError,
		// 	expectedMessage:   "failed to parse locations",
		// 	mockRedisSetup:    nil,
		// },
		// {
		// 	testName:          "large CSV",
		// 	sheetId:           "mockSheetId",
		// 	mockCsvResponse:   generateLargeCSV(),
		// 	mockCsvStatusCode: http.StatusOK,
		// 	expectedStatus:    http.StatusOK,
		// 	expectedMessage:   `<iframe src="/?sheetId=mockSheetId" width="100%" height="300" frameborder="0"></iframe>`,
		// 	mockRedisSetup: func(mockRedisClient *testUtils.MockRedisClient, ctx context.Context) {
		// 		mockRedisClient.On("Set",
		// 			mock.Anything,
		// 			mock.AnythingOfType("string"),
		// 			mock.AnythingOfType("[]uint8"),
		// 			mock.AnythingOfType("time.Duration"),
		// 		).Return(&redis.StatusCmd{})
		// 	},
		// },
		// {
		// 	testName:          "slow CSV server",
		// 	sheetId:           "mockSheetId",
		// 	mockCsvResponse:   "Name,Address,City,State,Country,Website,Phone Number,Latitude,Longitude\nLocation 1,Address 1,City 1,State 1,Country 1,http://example.com,1234567890,40.7128,-74.0060\n",
		// 	mockCsvStatusCode: http.StatusOK,
		// 	expectedStatus:    http.StatusInternalServerError,
		// 	expectedMessage:   "failed to fetch CSV data",
		// 	mockRedisSetup:    nil,
		// },
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
		runTestCase(t, testCase)
	}
}

func runTestCase(t *testing.T, testCase TestCaseConfig) {
	timeout := 2 * time.Second

	t.Run(testCase.testName, func(t *testing.T) {
		ctx, cancelCtxTimeout := context.WithTimeout(context.Background(), timeout)

		mockRedisClient := &testUtils.MockRedisClient{}

		mockCSVServer := createMockCSVServer(testCase)

		defer func() {
			cancelCtxTimeout()
			mockRedisClient.AssertExpectations(t)
			mockRedisClient.Calls = nil
			mockRedisClient.ExpectedCalls = nil
			mockCSVServer.Close()
		}()

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
	})
}

func generateLargeCSV() string {
	var b strings.Builder
	b.WriteString("Name,Address,City,State,Country,Website,Phone Number,Latitude,Longitude\n")
	for i := 0; i < 10000; i++ {
		fmt.Fprintf(&b, "Location %d,Address %d,City %d,State %d,Country %d,http://example.com,%d,%.6f,%.6f\n",
			i, i, i, i, i, i, rand.Float64()*90, rand.Float64()*180)
	}
	return b.String()
}

func createMockCSVServer(testCase TestCaseConfig) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if testCase.testName == "slow CSV server" {
			time.Sleep(3 * time.Second)
		}
		w.WriteHeader(testCase.mockCsvStatusCode)
		_, _ = w.Write([]byte(testCase.mockCsvResponse))
	}))
}
