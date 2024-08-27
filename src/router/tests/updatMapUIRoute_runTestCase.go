package router

import (
	"context"
	"gomap/src/router"
	"gomap/src/testUtils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestCaseConfig struct {
	testName          string
	sheetId           string
	mockCSVResponse   string
	mockCSVStatusCode int
	expectedStatus    int
	expectedMessage   string
	mockRedisSetup    func(mockRedisClient *testUtils.MockRedisClient, ctx context.Context)
	customAssert      func(t *testing.T, recorder *httptest.ResponseRecorder)
}

func runTestCase(t *testing.T, testCase TestCaseConfig) {
	timeout := 30 * time.Second
	testShouldHaveDelay := testCase.testName == "slow_CSV_server"
	// if testShouldHaveDelay {
	// 	timeout = 1 * time.Second
	// }

	t.Run(testCase.testName, func(t *testing.T) {
		ctx, cancelCtxTimeout := context.WithTimeout(context.Background(), timeout)

		mockRedisClient := &testUtils.MockRedisClient{}

		mockCSVServer := testUtils.CreateMockCSVServer(testUtils.MockCSVServerConfig{
			AddDelay:          testShouldHaveDelay,
			MockCSVResponse:   testCase.mockCSVResponse,
			MockCSVStatusCode: testCase.mockCSVStatusCode,
		})

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
		if testCase.customAssert != nil {
			testCase.customAssert(t, routerRecorder)
		}
	})
}
