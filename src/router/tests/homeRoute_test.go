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

var (
	demoPrompt      = "You are using the demo version."
	homepageTitle   = "MAPPERBOI BETA"
	homepageTagline = "Website Map Generator"
	mapIdString     = "id=\"map\""
)

func setupRouterTest(t *testing.T, queryParam string) (string, *html.Node) {
	ctx := context.Background()
	mockRedisClient := &testUtils.MockRedisClient{}
	router := router.InitRouter(mockRedisClient, ctx)

	req, err := http.NewRequest("GET", "/"+queryParam, nil)
	if err != nil {
		t.Fatal(err)
	}

	routerRecorder := httptest.NewRecorder()
	router.ServeHTTP(routerRecorder, req)

	if routerRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, routerRecorder.Code)
	}

	stringifiedDoc := routerRecorder.Body.String()
	htmlDoc, err := html.Parse(strings.NewReader(stringifiedDoc))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	return stringifiedDoc, htmlDoc
}

func TestHomeRouter(t *testing.T) {
	stringifiedDoc, htmlDoc := setupRouterTest(t, "")

	testUtils.CheckElement(t, htmlDoc, "h1", homepageTitle)

	if strings.Contains(stringifiedDoc, demoPrompt) {
		t.Errorf("Demo section displayed without demo=true query param")
	}

	if strings.Contains(stringifiedDoc, mapIdString) {
		t.Errorf("Homepage map state displayed without sheetId query param")
	}
}

func TestHomeRouter_WithDemo(t *testing.T) {
	stringifiedDoc, _ := setupRouterTest(t, "?demo=true")

	if !strings.Contains(stringifiedDoc, demoPrompt) {
		t.Errorf("Demo section NOT displayed with demo=true query param")
	}
}

func TestHomeRouter_WithSheetId(t *testing.T) {
	t.Skip("Skipping this test")
	stringifiedDoc, _ := setupRouterTest(t, "?sheetId=2PACK-3kj3l2kjf32f")

	if !strings.Contains(stringifiedDoc, mapIdString) {
		t.Errorf("Map state NOT displayed with sheetId query param")
	}
}
