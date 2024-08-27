package router

import (
	"context"
	"errors"
	"gomap/src/testUtils"
	"net/http"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

const defaultValidCSVResponse = "Name,Address,City,State,Country,Website,Phone Number,Latitude,Longitude\nLocation 1,Address 1,City 1,State 1,Country 1,http://example.com,1234567890,40.7128,-74.0060\n"

func TestUpdateMapUIHandlerIntegration(t *testing.T) {

	testCases := []TestCaseConfig{
		{
			testName:          "processed_locations_successfully",
			sheetId:           "mockSheetId",
			mockCSVResponse:   defaultValidCSVResponse,
			mockCSVStatusCode: http.StatusOK,
			expectedStatus:    http.StatusOK,
			expectedMessage:   `<iframe src="/?sheetId=mockSheetId" width="100%" height="300" frameborder="0"></iframe>`,
			mockRedisSetup: func(mockRedisClient *testUtils.MockRedisClient, ctx context.Context) {
				mockRedisClient.On("Set",
					mock.Anything,
					"mockSheetId",
					mock.AnythingOfType("[]uint8"),
					mock.AnythingOfType("time.Duration"),
				).Return(&redis.StatusCmd{})
			},
		},
		{
			testName:          "no_sheetId_provided",
			sheetId:           "",
			mockCSVResponse:   "",
			mockCSVStatusCode: http.StatusOK,
			expectedStatus:    http.StatusBadRequest,
			expectedMessage:   "Missing sheetId parameter",
			mockRedisSetup:    nil,
		},
		{
			testName:          "fetch_locations_fails",
			sheetId:           "mockSheetId",
			mockCSVResponse:   "",
			mockCSVStatusCode: http.StatusInternalServerError,
			expectedStatus:    http.StatusInternalServerError,
			expectedMessage:   "failed to fetch CSV data",
			mockRedisSetup:    nil,
		},
		{
			testName: "parse_locations_fails",
			sheetId:  "mockSheetId",
			mockCSVResponse: `Name,Address,City,State,Country,Website,Phone Number,Latitude,Longitude
		    Location 1,Address 1,City 1,State 1,Country 1,https://example.com,1234567890,INVALID_LAT,INVALID_LONG`,
			mockCSVStatusCode: http.StatusOK,
			expectedStatus:    http.StatusInternalServerError,
			expectedMessage:   "failed to parse locations",
			mockRedisSetup:    nil,
		},
		{
			testName:          "invalid_CSV_format",
			sheetId:           "mockSheetId",
			mockCSVResponse:   "Name,Address\nLocation 1,Address 1\n",
			mockCSVStatusCode: http.StatusOK,
			expectedStatus:    http.StatusInternalServerError,
			expectedMessage:   "failed to parse locations",
			mockRedisSetup:    nil,
		},
		{
			testName:          "empty_CSV",
			sheetId:           "mockSheetId",
			mockCSVResponse:   "",
			mockCSVStatusCode: http.StatusOK,
			expectedStatus:    http.StatusInternalServerError,
			expectedMessage:   "no valid locations found in CSV",
			mockRedisSetup:    nil,
		},
		// {
		// 	testName:          "slow_CSV_server",
		// 	sheetId:           "mockSheetId",
		// 	mockCSVResponse:   defaultValidCSVResponse,
		// 	mockCSVStatusCode: http.StatusOK,
		// 	expectedStatus:    http.StatusInternalServerError,
		// 	expectedMessage:   "the operation timed out, please check your spreadsheet and try again",
		// 	mockRedisSetup:    nil,
		// },
		{
			testName:          "process_locations_fails",
			sheetId:           "mockSheetId",
			mockCSVResponse:   defaultValidCSVResponse,
			mockCSVStatusCode: http.StatusOK,
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
