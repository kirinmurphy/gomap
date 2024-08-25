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

	"github.com/redis/go-redis/v9"
)

func TestGetLocationsRoute(t *testing.T) {
	ctx := context.Background()
	mockRedisClient := &testUtils.MockRedisClient{}

	r := router.InitRouter(mockRedisClient, ctx)

	tests := []struct {
		name            string
		sheetId         string
		setupMock       func()
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:    "Missing sheetId",
			sheetId: "",
			setupMock: func() {
				mockRedisClient.GetFunc = func(ctx context.Context, key string) *redis.StringCmd {
					return redis.NewStringResult("", redis.Nil)
				}
			},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Missing sheetId parameter",
		},
		{
			name:    "No locations found",
			sheetId: "nonexistent",
			setupMock: func() {
				mockRedisClient.GetFunc = func(ctx context.Context, key string) *redis.StringCmd {
					return redis.NewStringResult("", redis.Nil)
				}
			},
			expectedStatus:  http.StatusNotFound,
			expectedMessage: "No locations found",
		},
		{
			name:    "Redis Internal Error",
			sheetId: "errorSheet",
			setupMock: func() {
				mockRedisClient.GetFunc = func(ctx context.Context, key string) *redis.StringCmd {
					cmd := redis.NewStringResult("", nil)
					cmd.SetErr(errors.New("Redis error"))
					return cmd
				}
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Redis error",
		},
		{
			name:    "Valid locations found",
			sheetId: "validSheet",
			setupMock: func() {
				mockRedisClient.GetFunc = func(ctx context.Context, key string) *redis.StringCmd {
					return redis.NewStringResult(`{"location": "data"}`, nil)
				}
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: `{"location": "data"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setupMock()

			req, err := http.NewRequest("GET", "/getLocations?sheetId="+test.sheetId, nil)
			if err != nil {
				t.Fatal(err)
			}

			routerRecorder := httptest.NewRecorder()
			r.ServeHTTP(routerRecorder, req)

			if routerRecorder.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, got %d", test.expectedStatus, routerRecorder.Code)
			}

			if strings.TrimSpace(routerRecorder.Body.String()) != strings.TrimSpace(test.expectedMessage) {
				t.Errorf("Expected message %s, got %s", test.expectedMessage, routerRecorder.Body.String())
			}
		})
	}
}
