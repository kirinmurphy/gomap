package testUtils

import (
	"log"
	"net/http"
	"net/http/httptest"
	"time"
)

const MOCK_CSV_SERVER_DELAY_SECONDS = 10

type MockCSVServerConfig struct {
	AddDelay          bool
	MockCSVResponse   string
	MockCSVStatusCode int
}

func CreateMockCSVServer(mocks MockCSVServerConfig) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mocks.AddDelay {
			select {
			case <-time.After(MOCK_CSV_SERVER_DELAY_SECONDS * time.Second):
			case <-r.Context().Done():
				log.Println("Request context cancelled in slow CSV server")
				return
			}
		}
		w.WriteHeader(mocks.MockCSVStatusCode)
		_, _ = w.Write([]byte(mocks.MockCSVResponse))
	}))
}
