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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetLocationsRoute(t *testing.T) {
	contextMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })

	tests := []struct {
		name            string
		sheetId         string
		setupMock       func(mockRedisClient *testUtils.MockRedisClient)
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "Missing sheetId",
			sheetId:         "",
			setupMock:       func(mockRedisClient *testUtils.MockRedisClient) {},
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Missing sheetId parameter",
		},
		{
			name:    "No locations found",
			sheetId: "nonexistent",
			setupMock: func(mockRedisClient *testUtils.MockRedisClient) {
				mockRedisClient.On("Get", contextMatcher, "nonexistent").Return(
					redis.NewStringResult("", redis.Nil),
				).Once()
			},
			expectedStatus:  http.StatusNotFound,
			expectedMessage: "No locations found",
		},
		{
			name:    "Redis Internal Error",
			sheetId: "errorSheet",
			setupMock: func(mockRedisClient *testUtils.MockRedisClient) {
				cmd := redis.NewStringResult("", nil)
				cmd.SetErr(errors.New("Redis error"))
				mockRedisClient.On("Get", contextMatcher, "errorSheet").Return(cmd).Once()
			},
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "Redis error",
		},
		{
			name:    "Valid locations found",
			sheetId: "validSheet",
			setupMock: func(mockRedisClient *testUtils.MockRedisClient) {
				mockRedisClient.On("Get", contextMatcher, "validSheet").Return(
					redis.NewStringResult(`{"location": "data"}`, nil),
				).Once()
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: `{"location": "data"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			mockRedisClient := new(testUtils.MockRedisClient)

			test.setupMock(mockRedisClient)

			r := router.InitRouter(router.RouterConfig{
				RedisClient:        mockRedisClient,
				Ctx:                ctx,
				BaseSpreadsheetUrl: "asdf_%s_asdf",
			})

			req, err := http.NewRequest("GET", "/getLocations?sheetId="+test.sheetId, nil)
			if err != nil {
				t.Fatal(err)
			}

			routerRecorder := httptest.NewRecorder()
			r.ServeHTTP(routerRecorder, req)

			assert.Equal(t, test.expectedStatus, routerRecorder.Code)
			assert.Equal(t, strings.TrimSpace(test.expectedMessage), strings.TrimSpace(routerRecorder.Body.String()))

			mockRedisClient.AssertExpectations(t)
		})
	}
}
