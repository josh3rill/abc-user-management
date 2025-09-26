package services

import (
    "context"
    "time"
    "strings"
    
    "github.com/go-redis/redis/v8"
)

type CacheService struct {
    client *redis.Client
    ctx    context.Context
}

// NewCacheService creates a new cache service instance
func NewCacheService(client *redis.Client) *CacheService {
    return &CacheService{
        client: client,
        ctx:    context.Background(),
    }
}

// Set stores value in cache with expiration
func (c *CacheService) Set(key string, value string, expiration time.Duration) error {
    return c.client.Set(c.ctx, key, value, expiration).Err()
}

// Get retrieves value from cache
func (c *CacheService) Get(key string) (string, error) {
    val, err := c.client.Get(c.ctx, key).Result()
    if err == redis.Nil {
        return "", nil // Key doesn't exist
    }
    return val, err
}

// Delete removes keys from cache, supports wildcards
func (c *CacheService) Delete(pattern string) error {
    // Handle wildcard patterns for bulk deletion
    if strings.Contains(pattern, "*") {
        keys, err := c.client.Keys(c.ctx, pattern).Result()
        if err != nil {
            return err
        }
        if len(keys) > 0 {
            return c.client.Del(c.ctx, keys...).Err()
        }
        return nil
    }
    // Delete single key
    return c.client.Del(c.ctx, pattern).Err()
}

// Exists checks if key exists in cache
func (c *CacheService) Exists(key string) (bool, error) {
    val, err := c.client.Exists(c.ctx, key).Result()
    if err != nil {
        return false, err
    }
    return val > 0, nil
}

// SetNX sets value only if key doesn't exist
func (c *CacheService) SetNX(key string, value string, expiration time.Duration) (bool, error) {
    return c.client.SetNX(c.ctx, key, value, expiration).Result()
}

// IncrementCounter increments numeric value
func (c *CacheService) IncrementCounter(key string) (int64, error) {
    return c.client.Incr(c.ctx, key).Result()
}

// GetMultiple retrieves multiple keys at once
func (c *CacheService) GetMultiple(keys []string) (map[string]string, error) {
    result := make(map[string]string)
    
    // Use pipeline for efficiency
    pipe := c.client.Pipeline()
    cmds := make([]*redis.StringCmd, len(keys))
    
    for i, key := range keys {
        cmds[i] = pipe.Get(c.ctx, key)
    }
    
    _, err := pipe.Exec(c.ctx)
    if err != nil && err != redis.Nil {
        return nil, err
    }
    
    // Collect results
    for i, cmd := range cmds {
        val, err := cmd.Result()
        if err == nil {
            result[keys[i]] = val
        }
    }
    
    return result, nil
}