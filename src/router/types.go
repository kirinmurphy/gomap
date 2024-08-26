package router

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RouterConfig struct {
	RedisClient        RedisClientInterface
	Ctx                context.Context
	BaseSpreadsheetUrl string
}

type RedisClientInterface interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}
