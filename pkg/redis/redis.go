package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xsj/gin-sqlc/config"
	"github.com/0xsj/gin-sqlc/log"
	"github.com/go-redis/redis/v8"
)

var Nil = redis.Nil

type Client struct {
	rdb    *redis.Client
	logger log.Logger
}

func NewClient(cfg config.Config, logger log.Logger) (*Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	logger.Infof("Connecting to Redis at %s", addr)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Errorf("Failed to connect to Redis: %v", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	logger.Info("Successfully connected to Redis")
	return &Client{
		rdb:    rdb,
		logger: logger,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.logger.Debugf("Setting Redis key: %s", key)
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key from Redis
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	c.logger.Debugf("Getting Redis key: %s", key)
	return c.rdb.Get(ctx, key).Result()
}

// Delete removes keys from Redis
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	c.logger.Debugf("Deleting Redis keys: %v", keys)
	return c.rdb.Del(ctx, keys...).Err()
}

// SetJSON stores a JSON-serializable object in Redis
func (c *Client) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.logger.Debugf("Setting JSON in Redis key: %s", key)
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return c.rdb.Set(ctx, key, data, expiration).Err()
}

// GetJSON retrieves a JSON-serialized object from Redis
func (c *Client) GetJSON(ctx context.Context, key string, dest interface{}) error {
	c.logger.Debugf("Getting JSON from Redis key: %s", key)
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Exists checks if a key exists in Redis
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	c.logger.Debugf("Checking if Redis key exists: %s", key)
	res, err := c.rdb.Exists(ctx, key).Result()
	return res > 0, err
}

// Expire sets an expiration time on a key
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	c.logger.Debugf("Setting expiration for Redis key: %s", key)
	return c.rdb.Expire(ctx, key, expiration).Err()
}

// Incr increments the value of a key
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	c.logger.Debugf("Incrementing Redis key: %s", key)
	return c.rdb.Incr(ctx, key).Result()
}

// Keys returns all keys matching a pattern
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	c.logger.Debugf("Getting Redis keys matching pattern: %s", pattern)
	return c.rdb.Keys(ctx, pattern).Result()
}

// FlushDB removes all keys from the current database
func (c *Client) FlushDB(ctx context.Context) error {
	c.logger.Warn("Flushing Redis database")
	return c.rdb.FlushDB(ctx).Err()
}

// Close closes the Redis client connection
func (c *Client) Close() error {
	c.logger.Info("Closing Redis connection")
	return c.rdb.Close()
}
