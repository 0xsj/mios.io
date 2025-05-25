package cache

import (
	"context"
	"time"

	"github.com/0xsj/gin-sqlc/log"
	"github.com/0xsj/gin-sqlc/pkg/redis"
)

type CacheService interface {
	Get(ctx context.Context, key string, dest interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeletePattern(ctx context.Context, pattern string) error
	GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fetchFn func() (interface{}, error)) error
}

type RedisCache struct {
	client *redis.Client
	logger log.Logger
	prefix string
}
