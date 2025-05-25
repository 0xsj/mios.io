// pkg/cache/cache.go
package cache

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
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

func NewRedisCache(client *redis.Client, logger log.Logger, prefix string) CacheService {
	return &RedisCache{
		client: client,
		logger: logger,
		prefix: prefix,
	}
}

func (c *RedisCache) buildKey(key string) string {
	if c.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	fullKey := c.buildKey(key)
	
	data, err := c.client.Get(ctx, fullKey)
	if err != nil {
		if err == redis.Nil {
			return false, nil // Cache miss
		}
		c.logger.Errorf("Cache get error for key %s: %v", fullKey, err)
		return false, err
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		c.logger.Errorf("Cache unmarshal error for key %s: %v", fullKey, err)
		return false, err
	}

	c.logger.Debugf("Cache hit for key: %s", fullKey)
	return true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := c.buildKey(key)
	
	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Errorf("Cache marshal error for key %s: %v", fullKey, err)
		return err
	}

	if err := c.client.Set(ctx, fullKey, data, ttl); err != nil {
		c.logger.Errorf("Cache set error for key %s: %v", fullKey, err)
		return err
	}

	c.logger.Debugf("Cache set for key: %s (TTL: %v)", fullKey, ttl)
	return nil
}

func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	fullKeys := make([]string, len(keys))
	for i, key := range keys {
		fullKeys[i] = c.buildKey(key)
	}

	if err := c.client.Delete(ctx, fullKeys...); err != nil {
		c.logger.Errorf("Cache delete error for keys %v: %v", fullKeys, err)
		return err
	}

	c.logger.Debugf("Cache deleted keys: %v", fullKeys)
	return nil
}

func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	fullPattern := c.buildKey(pattern)
	
	keys, err := c.client.Keys(ctx, fullPattern)
	if err != nil {
		c.logger.Errorf("Cache keys scan error for pattern %s: %v", fullPattern, err)
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	if err := c.client.Delete(ctx, keys...); err != nil {
		c.logger.Errorf("Cache delete pattern error for pattern %s: %v", fullPattern, err)
		return err
	}

	c.logger.Debugf("Cache deleted %d keys matching pattern: %s", len(keys), fullPattern)
	return nil
}

func (c *RedisCache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fetchFn func() (interface{}, error)) error {
	// Try to get from cache first
	found, err := c.Get(ctx, key, dest)
	if err != nil {
		c.logger.Warnf("Cache get error, fetching fresh data: %v", err)
	}
	
	if found && err == nil {
		return nil // Cache hit
	}

	// Cache miss, fetch fresh data
	c.logger.Debugf("Cache miss for key: %s, fetching fresh data", key)
	
	value, err := fetchFn()
	if err != nil {
		return err
	}

	// Store in cache for next time
	if setErr := c.Set(ctx, key, value, ttl); setErr != nil {
		c.logger.Warnf("Failed to cache data for key %s: %v", key, setErr)
		// Don't return error, just log it since we have the data
	}

	// Marshal the fresh data into dest
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Cache key builders for different operations
type CacheKeyBuilder struct{}

func NewCacheKeyBuilder() *CacheKeyBuilder {
	return &CacheKeyBuilder{}
}

func (kb *CacheKeyBuilder) UserAnalytics(userID string, days int) string {
	return fmt.Sprintf("analytics:user:%s:days:%d", userID, days)
}

func (kb *CacheKeyBuilder) ProfileDashboard(userID string, days int) string {
	return fmt.Sprintf("dashboard:user:%s:days:%d", userID, days)
}

func (kb *CacheKeyBuilder) ContentItemAnalytics(itemID string, timeRange string) string {
	hash := kb.HashString(timeRange)
	return fmt.Sprintf("analytics:item:%s:range:%s", itemID, hash)
}

func (kb *CacheKeyBuilder) TimeRangeAnalytics(userID, startDate, endDate string) string {
	hash := kb.HashString(startDate + endDate)
	return fmt.Sprintf("analytics:user:%s:timerange:%s", userID, hash)
}

func (kb *CacheKeyBuilder) ReferrerAnalytics(userID, startDate, endDate string, limit int) string {
	hash := kb.HashString(fmt.Sprintf("%s:%s:%d", startDate, endDate, limit))
	return fmt.Sprintf("analytics:referrer:%s:range:%s", userID, hash)
}

func (kb *CacheKeyBuilder) PageViewAnalytics(userID, startDate, endDate string, limit int) string {
	hash := kb.HashString(fmt.Sprintf("%s:%s:%d", startDate, endDate, limit))
	return fmt.Sprintf("analytics:pageviews:user:%s:range:%s", userID, hash)
}

func (kb *CacheKeyBuilder) LinkMetadata(url string) string {
	hash := kb.HashString(url)
	return fmt.Sprintf("metadata:url:%s", hash)
}

func (kb *CacheKeyBuilder) UserContent(userID string) string {
	return fmt.Sprintf("content:user:%s", userID)
}

// HashString is exported so it can be used from other packages
func (kb *CacheKeyBuilder) HashString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))[:8]
}

// Invalidation patterns
func (kb *CacheKeyBuilder) UserAnalyticsPattern(userID string) string {
	return fmt.Sprintf("analytics:user:%s:*", userID)
}

func (kb *CacheKeyBuilder) UserDashboardPattern(userID string) string {
	return fmt.Sprintf("dashboard:user:%s:*", userID)
}

func (kb *CacheKeyBuilder) UserContentPattern(userID string) string {
	return fmt.Sprintf("content:user:%s*", userID)
}

// Cache TTL constants
const (
	DefaultTTL = 5 * time.Minute
	ShortTTL   = 1 * time.Minute
	MediumTTL  = 15 * time.Minute
	LongTTL    = 1 * time.Hour
	DayTTL     = 24 * time.Hour
)

// TTL strategies for different data types
func GetAnalyticsTTL() time.Duration {
	return MediumTTL // Analytics can be slightly stale
}

func GetDashboardTTL() time.Duration {
	return ShortTTL // Dashboard should be relatively fresh
}

func GetContentTTL() time.Duration {
	return DefaultTTL // Content changes require quick updates
}

func GetMetadataTTL() time.Duration {
	return DayTTL // Link metadata rarely changes
}

func GetExpensiveOperationTTL() time.Duration {
	return LongTTL // Long TTL for expensive operations
}