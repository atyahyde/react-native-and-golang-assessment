package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// listKeysTrackingSet adalah Redis SET yang menyimpan semua cache key list
// task yang sedang aktif, supaya invalidasi tidak perlu memakai KEYS/SCAN
// (yang mahal di Redis produksi) dan cukup DEL anggota-anggota set ini.
const listKeysTrackingSet = "tasks:list:keys"

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) TrackListKey(ctx context.Context, key string) error {
	return c.client.SAdd(ctx, listKeysTrackingSet, key).Err()
}

func (c *RedisCache) InvalidateListCache(ctx context.Context) error {
	keys, err := c.client.SMembers(ctx, listKeysTrackingSet).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}
	return c.client.Del(ctx, listKeysTrackingSet).Err()
}
