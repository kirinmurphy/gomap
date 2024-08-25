package router

import (
	"context"
	"fmt"
	"gomap/src/router"
	"gomap/src/testUtils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/html"
)

var (
	demoPrompt    = "You are using the demo version."
	homepageTitle = "MAPPERBOI BETA"
	mapIdString   = "id=\"map\""
)

func setupRouterTest(t *testing.T, queryParam string, redisGetFunc func(ctx context.Context, key string) *redis.StringCmd) (string, *html.Node) {
	ctx := context.Background()

	mockRedisClient := &testUtils.MockRedisClient{}
	if redisGetFunc != nil {
		mockRedisClient.GetFunc = redisGetFunc
	}

	r := router.InitRouter(mockRedisClient, ctx)

	req, err := http.NewRequest("GET", "/"+queryParam, nil)
	if err != nil {
		t.Fatal(err)
	}

	routerRecorder := httptest.NewRecorder()
	r.ServeHTTP(routerRecorder, req)

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

func TestHomeRouter_Default(t *testing.T) {
	stringifiedDoc, htmlDoc := setupRouterTest(t, "", nil)

	testUtils.CheckElement(t, htmlDoc, "h1", homepageTitle)

	if strings.Contains(stringifiedDoc, demoPrompt) {
		t.Errorf("Demo section displayed without demo=true query param")
	}

	if strings.Contains(stringifiedDoc, mapIdString) {
		t.Errorf("Homepage map state displayed without sheetId query param")
	}
}

func TestHomeRouter_WithDemo(t *testing.T) {
	stringifiedDoc, _ := setupRouterTest(t, "?demo=true", nil)

	if !strings.Contains(stringifiedDoc, demoPrompt) {
		t.Errorf("Demo section NOT displayed with demo=true query param")
	}
}

func TestHomeRouter_WithSheetId(t *testing.T) {
	sheetId := "2PACK-3kj3l2kjf32f"
	sheetIdParam := "?sheetId=" + sheetId

	stringifiedDoc, _ := setupRouterTest(t, sheetIdParam,
		func(ctx context.Context, key string) *redis.StringCmd {
			if key == sheetId {
				return redis.NewStringResult("mocked result", nil)
			}
			return redis.NewStringResult("", redis.Nil)
		},
	)

	if !strings.Contains(stringifiedDoc, mapIdString) {
		t.Errorf("Map state NOT displayed with sheetId query param")
	}

	hxMapQuery := fmt.Sprintf(`hx-get="/getLocations?sheetId=%s"`, sheetId)
	if !strings.Contains(stringifiedDoc, hxMapQuery) {
		t.Errorf("hxQuery for map data NOT present")
	}
}
