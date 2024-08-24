package router

import (
	"context"
	"gomap/src/router"
	"gomap/src/testUtils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestHomeRouter(t *testing.T) {
	ctx := context.Background()
	mockRedisClient := &testUtils.MockRedisClient{}
	router := router.InitRouter(mockRedisClient, ctx)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	routerRecorder := httptest.NewRecorder()
	router.ServeHTTP(routerRecorder, req)

	if routerRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, routerRecorder.Code)
	}

	doc, err := html.Parse(strings.NewReader(routerRecorder.Body.String()))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	var sb strings.Builder
	err = html.Render(&sb, doc)
	if err != nil {
		t.Fatalf("Failed to render HTML: %v", err)
	}

	testUtils.CheckElement(t, doc, "h1", "MAPPERBOI BETA")
}
