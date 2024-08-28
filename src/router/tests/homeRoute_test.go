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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/html"
)

var (
	demoPrompt    = "You are using the demo version."
	homepageTitle = "MAPPERIZER BETA"
	mapIdString   = "id=\"map\""
)

func TestHomeRouter(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		stringifiedDoc, htmlDoc := setupRouterTest(RouterTestConfig{
			T:          t,
			QueryParam: "",
		})

		testUtils.CheckElement(t, htmlDoc, "h1", homepageTitle)

		assert.NotContains(t, stringifiedDoc, demoPrompt, "Demo section displayed without demo=true query param")
		assert.NotContains(t, stringifiedDoc, mapIdString, "Homepage map state displayed without sheetId query param")
	})

	t.Run("WithDemo", func(t *testing.T) {
		stringifiedDoc, _ := setupRouterTest(RouterTestConfig{
			T:          t,
			QueryParam: "?demo=true",
		})
		assert.Contains(t, stringifiedDoc, demoPrompt, "Demo section NOT displayed with demo=true query param")
	})

	t.Run("WithSheetId", func(t *testing.T) {
		sheetId := "2PACK-3kj3l2kjf32f"
		sheetIdParam := "?sheetId=" + sheetId

		mockRedisClient := &testUtils.MockRedisClient{}
		contextMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
		keyMatcher := mock.MatchedBy(func(key string) bool { return key == sheetId })

		mockRedisClient.On("Get", contextMatcher, keyMatcher).Return(
			redis.NewStringResult(("mocked result"), nil),
		)

		stringifiedDoc, _ := setupRouterTest(RouterTestConfig{
			T:               t,
			QueryParam:      sheetIdParam,
			MockRedisClient: mockRedisClient,
		})

		assert.Contains(t, stringifiedDoc, mapIdString, "Map state NOT displayed with sheetId query param")

		hxMapQuery := fmt.Sprintf(`hx-get="/getLocations?sheetId=%s"`, sheetId)
		assert.Contains(t, stringifiedDoc, hxMapQuery, "hxQuery for map data NOT present")
	})
}

type RouterTestConfig struct {
	T               *testing.T
	QueryParam      string
	MockRedisClient *testUtils.MockRedisClient
}

func setupRouterTest(testConfig RouterTestConfig) (string, *html.Node) {
	ctx := context.Background()

	if testConfig.MockRedisClient == nil {
		testConfig.MockRedisClient = &testUtils.MockRedisClient{}
	}

	r := router.InitRouter(router.RouterConfig{
		RedisClient:        testConfig.MockRedisClient,
		Ctx:                ctx,
		BaseSpreadsheetUrl: "asdf_%s_asdf",
	})

	req, err := http.NewRequest("GET", "/"+testConfig.QueryParam, nil)
	if err != nil {
		testConfig.T.Fatal(err)
	}

	routerRecorder := httptest.NewRecorder()
	r.ServeHTTP(routerRecorder, req)

	assert.Equal(testConfig.T, http.StatusOK, routerRecorder.Code)

	stringifiedDoc := routerRecorder.Body.String()
	htmlDoc, err := html.Parse(strings.NewReader(stringifiedDoc))
	assert.NoError(testConfig.T, err)

	return stringifiedDoc, htmlDoc
}
