package cache

import (
	"context"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	redisConfig := configuration.GetRedisConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisConfig.URL,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		PoolSize:     redisConfig.Pool.Max,                                   // max total connections
		MinIdleConns: redisConfig.Pool.Idle,                                  // minimum idle connections
		PoolTimeout:  time.Duration(redisConfig.Pool.Lifetime) * time.Second, // pool timeout
	})

	return &RedisCache{
		client: rdb,
	}
}

// Set sets a key-value pair in the cache with expiration
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value from the cache by key
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Delete removes a key from the cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in the cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.client.Exists(ctx, key).Result()
	return exists > 0, err
}

// Close closes the Redis client connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// GetClient returns the underlying Redis client
func (c *RedisCache) GetClient() *redis.Client {
	return c.client
}
