// pkg/cache/instrumented_cache.go
package cache

import (
	"context"
	"time"

	"github.com/0xsj/mios.io/pkg/metrics"
)

// InstrumentedCache wraps CacheService with metrics
type InstrumentedCache struct {
	base    CacheService
	metrics *metrics.Metrics
	name    string // cache instance name for metrics
}

func NewInstrumentedCache(base CacheService, metrics *metrics.Metrics, name string) CacheService {
	return &InstrumentedCache{
		base:    base,
		metrics: metrics,
		name:    name,
	}
}

func (c *InstrumentedCache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	start := time.Now()
	found, err := c.base.Get(ctx, key, dest)
	duration := time.Since(start)

	result := "miss"
	if found && err == nil {
		result = "hit"
	} else if err != nil {
		result = "error"
	}

	c.metrics.RecordCacheOperation("get", result, duration)
	
	// Update hit ratio
	c.updateHitRatio(result == "hit")

	return found, err
}

func (c *InstrumentedCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	start := time.Now()
	err := c.base.Set(ctx, key, value, ttl)
	duration := time.Since(start)

	result := "success"
	if err != nil {
		result = "error"
	}

	c.metrics.RecordCacheOperation("set", result, duration)
	return err
}

func (c *InstrumentedCache) Delete(ctx context.Context, keys ...string) error {
	start := time.Now()
	err := c.base.Delete(ctx, keys...)
	duration := time.Since(start)

	result := "success"
	if err != nil {
		result = "error"
	}

	c.metrics.RecordCacheOperation("delete", result, duration)
	return err
}

func (c *InstrumentedCache) DeletePattern(ctx context.Context, pattern string) error {
	start := time.Now()
	err := c.base.DeletePattern(ctx, pattern)
	duration := time.Since(start)

	result := "success"
	if err != nil {
		result = "error"
	}

	c.metrics.RecordCacheOperation("delete_pattern", result, duration)
	return err
}

func (c *InstrumentedCache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, fetchFn func() (interface{}, error)) error {
	start := time.Now()
	
	// First try to get from cache
	found, err := c.base.Get(ctx, key, dest)
	if found && err == nil {
		// Cache hit
		c.metrics.RecordCacheOperation("get", "hit", time.Since(start))
		c.updateHitRatio(true)
		return nil
	}

	// Cache miss, need to fetch and set
	value, fetchErr := fetchFn()
	if fetchErr != nil {
		c.metrics.RecordCacheOperation("get", "miss", time.Since(start))
		c.updateHitRatio(false)
		return fetchErr
	}

	// Set in cache
	if setErr := c.base.Set(ctx, key, value, ttl); setErr != nil {
		c.metrics.RecordCacheOperation("set", "error", time.Since(start))
	} else {
		c.metrics.RecordCacheOperation("set", "success", time.Since(start))
	}

	// Marshal the data into dest
	if data, err := jsonMarshal(value); err == nil {
		if err := jsonUnmarshal(data, dest); err == nil {
			c.metrics.RecordCacheOperation("get", "miss", time.Since(start))
			c.updateHitRatio(false)
			return nil
		}
	}

	c.metrics.RecordCacheOperation("get", "error", time.Since(start))
	c.updateHitRatio(false)
	return err
}

// Simple hit ratio calculation (you might want to use a more sophisticated sliding window)
func (c *InstrumentedCache) updateHitRatio(isHit bool) {
	// This is a simple implementation - in production you'd want to use
	// a sliding window or exponential moving average
	if isHit {
		c.metrics.CacheHitRatio.WithLabelValues(c.name).Set(1.0)
	} else {
		c.metrics.CacheHitRatio.WithLabelValues(c.name).Set(0.0)
	}
}

// Helper functions (you'd need to import encoding/json)
func jsonMarshal(v interface{}) ([]byte, error) {
	// Implementation would use json.Marshal
	return nil, nil
}

func jsonUnmarshal(data []byte, v interface{}) error {
	// Implementation would use json.Unmarshal
	return nil
}