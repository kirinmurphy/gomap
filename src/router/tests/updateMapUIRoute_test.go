package router

// import (
//   "context"
//   "errors"
//   "gomap/src/router"
//   "gomap/src/testUtils"
//   "net/http"
//   "net/http/httptest"
//   "strings"
//   "testing"

//   "github.com/stretchr/testify/assert"
//   "github.com/stretchr/testify/mock"
// )

// func TestUpdateMapUIHandlerIntegration(t *testing.T) {
//   ctx := context.Background()
//   mockRedisClient := &testUtils.MockRedisClient{}

//   mockCSVServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//     w.WriteHeader(http.StatusOK)
//     _, _ = w.Write([]byte("Name,Address,City,State,Country,Website,Phone Number,Latitude,Longitude\n" +
//       "Location 1,Address 1,City 1,State 1,Country 1,https://example.com, 1234567890,10.1234,20.1234\n"))
//   }))
//   defer mockCSVServer.Close()

//   // Replace the baseSpreadsheetUrl to use our mock server
//   router.SetBaseSpreadsheetUrl(mockCSVServer.URL + "?sheetId=%s")

//   r := router.InitRouter(mockRedisClient, ctx)

//   t.Run("no sheetId provided", func(t *testing.T) {
//     req, err := http.NewRequest("POST", "/updateMapUI", nil)
//     if err != nil {
//       t.Fatal(err)
//     }

//     req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//     routerRecorder := httptest.NewRecorder()
//     r.ServeHTTP(routerRecorder, req)

//     assert.Equal(t, http.StatusBadRequest, routerRecorder.Code)
//     assert.Contains(t, routerRecorder.Body.String(), "Missing sheetId parameter")
//   })

//   t.Run("fetch and parse locations successfully", func(t *testing.T) {
//     mockRedisClient.On("Set", ctx, "mockSheetId", mock.Anything, 0).Return(nil)

//     form := strings.NewReader("sheetId=mockSheetId")
//     req, err := http.NewRequest("POST", "/updateMapUI", form)
//     if err != nil {
//       t.Fatal(err)
//     }

//     req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//     routerRecorder := httptest.NewRecorder()
//     r.ServeHTTP(routerRecorder, req)

//     assert.Equal(t, http.StatusOK, routerRecorder.Code)
//     assert.Contains(t, routerRecorder.Body.String(), "mockSheetId")
//   })

//   t.Run("fetch locations fails", func(t *testing.T) {
//     mockCSVServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//       w.WriteHeader(http.StatusInternalServerError)
//     })

//     form := strings.NewReader("sheetId=mockSheetId")
//     req, err := http.NewRequest("POST", "/updateMapUI", form)
//     if err != nil {
//       t.Fatal(err)
//     }

//     req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//     routerRecorder := httptest.NewRecorder()
//     r.ServeHTTP(routerRecorder, req)

//     assert.Equal(t, http.StatusInternalServerError, routerRecorder.Code)
//     assert.Contains(t, routerRecorder.Body.String(), "failed to fetch CSV data")
//   })

//   t.Run("parse locations fails", func(t *testing.T) {
//     mockCSVServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//       w.WriteHeader(http.StatusOK)
//       _, _ = w.Write([]byte("Invalid CSV Content"))
//     })

//     form := strings.NewReader("sheetId=mockSheetId")
//     req, err := http.NewRequest("POST", "/updateMapUI", form)
//     if err != nil {
//       t.Fatal(err)
//     }

//     req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//     routerRecorder := httptest.NewRecorder()
//     r.ServeHTTP(routerRecorder, req)

//     assert.Equal(t, http.StatusInternalServerError, routerRecorder.Code)
//     assert.Contains(t, routerRecorder.Body.String(), "failed to parse locations")
//   })

//   t.Run("process locations fails", func(t *testing.T) {
//     mockRedisClient.On("Set", ctx, "mockSheetId", mock.Anything, 0).Return(errors.New("failed to cache locations"))

//     mockCSVServer.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//       w.WriteHeader(http.StatusOK)
//       _, _ = w.Write([]byte("Name,Address,City,State,Country,Website,Phone Number,Latitude,Longitude\n" +
//         "Location 1,Address 1,City 1,State 1,Country 1,https://example.com,1234567890,10.1234,20.1234\n"))
//     })

//     form := strings.NewReader("sheetId=mockSheetId")
//     req, err := http.NewRequest("POST", "/updateMapUI", form)
//     if err != nil {
//       t.Fatal(err)
//     }

//     req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//     routerRecorder := httptest.NewRecorder()
//     r.ServeHTTP(routerRecorder, req)

//     assert.Equal(t, http.StatusInternalServerError, routerRecorder.Code)
//     assert.Contains(t, routerRecorder.Body.String(), "failed to cache locations")
//   })
// }
