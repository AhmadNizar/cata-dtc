package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
	prefix string
}

func NewCacheRepository(client *redis.Client, prefix string) *CacheRepository {
	return &CacheRepository{
		client: client,
		prefix: prefix,
	}
}

func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshaling value: %w", err)
	}

	fullKey := r.getKey(key)
	if err := r.client.Set(ctx, fullKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("setting cache: %w", err)
	}

	return nil
}

func (r *CacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := r.getKey(key)
	data, err := r.client.Get(ctx, fullKey).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("getting cache: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("unmarshaling cache data: %w", err)
	}

	return nil
}

func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	fullKey := r.getKey(key)
	if err := r.client.Del(ctx, fullKey).Err(); err != nil {
		return fmt.Errorf("deleting cache: %w", err)
	}
	return nil
}

func (r *CacheRepository) DeleteByPattern(ctx context.Context, pattern string) error {
	fullPattern := r.getKey(pattern)
	keys, err := r.client.Keys(ctx, fullPattern).Result()
	if err != nil {
		return fmt.Errorf("getting keys by pattern: %w", err)
	}

	if len(keys) == 0 {
		return nil
	}

	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("deleting keys: %w", err)
	}

	return nil
}

func (r *CacheRepository) getKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

var ErrCacheMiss = fmt.Errorf("cache miss")