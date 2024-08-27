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

	t.Run("no sheetId provided", func(t *testing.T) {
		// t.Skip("Skipping this test temporarily")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		mockRedisClient := &testUtils.MockRedisClient{}
		mockCSVServer := createMockCSVServer("", http.StatusOK)
		defer mockCSVServer.Close()

		r := router.InitRouter(router.RouterConfig{
			RedisClient:        mockRedisClient,
			Ctx:                ctx,
			BaseSpreadsheetUrl: mockCSVServer.URL + "?sheetId=%s",
		})

		req, err := http.NewRequest("POST", "/updateMapUI", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		routerRecorder := httptest.NewRecorder()
		r.ServeHTTP(routerRecorder, req)

		assert.Equal(t, http.StatusBadRequest, routerRecorder.Code)
		assert.Contains(t, routerRecorder.Body.String(), "Missing sheetId parameter")

		defer func() {
			mockRedisClient.AssertExpectations(t)
			mockRedisClient.Calls = nil
			mockRedisClient.ExpectedCalls = nil
		}()
	})

	t.Run("fetch and parse locations successfully", func(t *testing.T) {
		// t.Skip("Skipping this test temporarily")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		mockRedisClient := &testUtils.MockRedisClient{}
		mockCSVServer := createMockCSVServer("", http.StatusOK)
		defer mockCSVServer.Close()

		r := router.InitRouter(router.RouterConfig{
			RedisClient:        mockRedisClient,
			Ctx:                ctx,
			BaseSpreadsheetUrl: mockCSVServer.URL + "?sheetId=%s",
		})

		mockRedisClient.On("Set",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]uint8"),
			mock.AnythingOfType("time.Duration"),
		).Return(&redis.StatusCmd{})

		form := strings.NewReader("sheetId=mockSheetId")
		req, err := http.NewRequest("POST", "/updateMapUI", form)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		routerRecorder := httptest.NewRecorder()
		r.ServeHTTP(routerRecorder, req)

		assert.Equal(t, http.StatusOK, routerRecorder.Code)
		assert.Contains(t, routerRecorder.Body.String(), "mockSheetId")

		defer func() {
			mockRedisClient.AssertExpectations(t)
			mockRedisClient.Calls = nil
			mockRedisClient.ExpectedCalls = nil
		}()
	})

	t.Run("fetch locations fails", func(t *testing.T) {
		// t.Skip("Skipping this test temporarily")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		mockRedisClient := &testUtils.MockRedisClient{}
		mockCSVServer := createMockCSVServer("", http.StatusInternalServerError)
		defer mockCSVServer.Close()

		r := router.InitRouter(router.RouterConfig{
			RedisClient:        mockRedisClient,
			Ctx:                ctx,
			BaseSpreadsheetUrl: mockCSVServer.URL + "?sheetId=%s",
		})

		form := strings.NewReader("sheetId=mockSheetId")
		req, err := http.NewRequest("POST", "/updateMapUI", form)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		routerRecorder := httptest.NewRecorder()
		r.ServeHTTP(routerRecorder, req)

		assert.Equal(t, http.StatusInternalServerError, routerRecorder.Code)
		assert.Contains(t, routerRecorder.Body.String(), "failed to fetch CSV data")

		mockRedisClient.AssertNotCalled(t, "Set")

		defer func() {
			mockRedisClient.AssertExpectations(t)
			mockRedisClient.Calls = nil
			mockRedisClient.ExpectedCalls = nil
		}()
	})

	t.Run("parse locations fails", func(t *testing.T) {
		// t.Skip("Skipping this test temporarily")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		mockRedisClient := &testUtils.MockRedisClient{}

		invalidCSV := "Name,Address,City,State,Country,Website,Phone Number,Latitude,Longitude\n" +
			"Location 1,Address 1,City 1,State 1,Country 1,https://example.com, 1234567890,INVALID_LAT,INVALID_LONG\n"

		mockCSVServer := createMockCSVServer(invalidCSV, http.StatusOK)
		defer mockCSVServer.Close()

		r := router.InitRouter(router.RouterConfig{
			RedisClient:        mockRedisClient,
			Ctx:                ctx,
			BaseSpreadsheetUrl: mockCSVServer.URL + "?sheetId=%s",
		})

		mockRedisClient.On("Set",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&redis.StatusCmd{})

		form := strings.NewReader("sheetId=mockSheetId")
		req, err := http.NewRequest("POST", "/updateMapUI", form)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		routerRecorder := httptest.NewRecorder()
		r.ServeHTTP(routerRecorder, req)

		assert.Equal(t, http.StatusInternalServerError, routerRecorder.Code)
		assert.Contains(t, routerRecorder.Body.String(), "failed to parse locations")

		defer func() {
			mockRedisClient.AssertExpectations(t)
			mockRedisClient.Calls = nil
			mockRedisClient.ExpectedCalls = nil
		}()
	})

	t.Run("process locations fails", func(t *testing.T) {
		t.Skip("Skipping this test temporarily")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		mockRedisClient := &testUtils.MockRedisClient{}
		mockRedisClient.AssertExpectations(t)
		mockRedisClient.ExpectedCalls = nil

		mockCSVServer := createMockCSVServer("", http.StatusOK)
		defer mockCSVServer.Close()

		r := router.InitRouter(router.RouterConfig{
			RedisClient:        mockRedisClient,
			Ctx:                ctx,
			BaseSpreadsheetUrl: mockCSVServer.URL + "?sheetId=%s",
		})

		statusCmd := redis.NewStatusCmd(ctx)
		statusCmd.SetErr(errors.New("failed to cache locations"))
		mockRedisClient.On("Set",
			mock.Anything,
			"mockSheetId",
			mock.Anything,
			mock.Anything,
		).Return(statusCmd)

		form := strings.NewReader("sheetId=mockSheetId")
		req, err := http.NewRequest("POST", "/updateMapUI", form)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		routerRecorder := httptest.NewRecorder()
		r.ServeHTTP(routerRecorder, req)

		assert.Equal(t, http.StatusInternalServerError, routerRecorder.Code)
		assert.Contains(t, routerRecorder.Body.String(), "failed to cache locations")

		defer func() {
			mockRedisClient.AssertExpectations(t)
			mockRedisClient.Calls = nil
			mockRedisClient.ExpectedCalls = nil
		}()
	})
}

func createMockCSVServer(response string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))
}
