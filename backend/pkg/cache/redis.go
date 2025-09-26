package cache

import (
    "context"
    "strings"
    "time"
    "abc-user-management/internal/config"
    
    "github.com/go-redis/redis/v8"
)

type RedisClient struct {
    Client *redis.Client
    ctx    context.Context
}

// InitRedis establishes connection to Redis cache
func InitRedis(cfg *config.Config) *RedisClient {
    rdb := redis.NewClient(&redis.Options{
        Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
        Password: cfg.RedisPassword,
        DB:       0,
    })
    
    return &RedisClient{
        Client: rdb,
        ctx:    context.Background(),
    }
}

// Set stores value in cache with expiration
func (r *RedisClient) Set(key string, value string, expiration time.Duration) error {
    return r.Client.Set(r.ctx, key, value, expiration).Err()
}

// Get retrieves value from cache
func (r *RedisClient) Get(key string) (string, error) {
    return r.Client.Get(r.ctx, key).Result()
}

// Delete removes key from cache
func (r *RedisClient) Delete(pattern string) error {
    // Support wildcard deletion for cache invalidation
    if strings.Contains(pattern, "*") {
        keys, err := r.Client.Keys(r.ctx, pattern).Result()
        if err != nil {
            return err
        }
        if len(keys) > 0 {
            return r.Client.Del(r.ctx, keys...).Err()
        }
        return nil
    }
    return r.Client.Del(r.ctx, pattern).Err()
}

// Close terminates Redis connection
func (r *RedisClient) Close() error {
    return r.Client.Close()
}