package router

// import (
// 	"context"
// 	"fmt"
// 	"gomap/src/testUtils"
// 	"math/rand/v2"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/redis/go-redis/v9"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestUpdateMapUIHandlerIntegration_LongTest(t *testing.T) {

// 	testCases := []TestCaseConfig{
// 		{
// 			testName:          "large_CSV",
// 			sheetId:           "mockSheetId",
// 			mockCSVResponse:   generateLargeCSV(),
// 			mockCSVStatusCode: http.StatusOK,
// 			expectedStatus:    http.StatusOK,
// 			expectedMessage:   `<iframe src="/?sheetId=mockSheetId" width="100%" height="300" frameborder="0"></iframe>`,
// 			mockRedisSetup: func(mockRedisClient *testUtils.MockRedisClient, ctx context.Context) {
// 				mockRedisClient.On("Set",
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 					mock.AnythingOfType("[]uint8"),
// 					mock.AnythingOfType("time.Duration"),
// 				).Return(&redis.StatusCmd{})
// 			},
// 			customAssert: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				t.Log("!!!!!!____" + recorder.Body.String())
// 				assert.Greater(t, len(recorder.Body.Bytes()), 1000000, "Expected response body to be large")
// 			},
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		runTestCase(t, testCase)
// 	}
// }

// func generateLargeCSV() string {
// 	var b strings.Builder
// 	b.WriteString("Name,Address,City,State,Country,Website,Phone Number,Latitude,Longitude\n")
// 	for i := 0; i < 10000; i++ {
// 		fmt.Fprintf(&b, "Location %d,Address %d,City %d,State %d,Country %d,http://example.com,%d,%.6f,%.6f\n",
// 			i, i, i, i, i, i, rand.Float64()*90, rand.Float64()*180)
// 	}
// 	return b.String()
// }
