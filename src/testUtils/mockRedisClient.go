package testUtils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type MockRedisClient struct{}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return redis.NewStringResult("", redis.Nil)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusResult("OK", nil)
}
