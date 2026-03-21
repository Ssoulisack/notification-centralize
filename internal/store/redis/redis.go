package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/your-org/notification-center/internal/config"
	"github.com/your-org/notification-center/internal/model"
)

type Cache struct {
	client *redis.Client
}

// Connect creates a Redis client.
func Connect(cfg config.RedisConfig) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Cache{client: client}, nil
}

// --- Rate Limiter ---

// AllowSend checks if a provider is within its rate limit.
// Uses a sliding window counter per provider.
func (c *Cache) AllowSend(ctx context.Context, providerName string, limit int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s", providerName)
	now := time.Now().UnixMilli()

	pipe := c.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-window.Milliseconds()))
	pipe.ZCard(ctx, key)
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
	pipe.Expire(ctx, key, window)

	results, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("rate limit check: %w", err)
	}

	count := results[1].(*redis.IntCmd).Val()
	return count < int64(limit), nil
}

// --- Device Token Cache ---

// CacheDeviceTokens stores device tokens in Redis for fast lookup.
func (c *Cache) CacheDeviceTokens(ctx context.Context, userID string, devices []*model.DeviceToken) error {
	key := fmt.Sprintf("devices:%s", userID)
	data, err := json.Marshal(devices)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, 10*time.Minute).Err()
}

// GetCachedDeviceTokens retrieves cached device tokens.
// Returns nil, nil if cache miss.
func (c *Cache) GetCachedDeviceTokens(ctx context.Context, userID string) ([]*model.DeviceToken, error) {
	key := fmt.Sprintf("devices:%s", userID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // cache miss
	}
	if err != nil {
		return nil, err
	}

	var devices []*model.DeviceToken
	if err := json.Unmarshal(data, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

// InvalidateDeviceCache clears cached tokens when a device registers or is removed.
func (c *Cache) InvalidateDeviceCache(ctx context.Context, userID string) error {
	return c.client.Del(ctx, fmt.Sprintf("devices:%s", userID)).Err()
}

// --- Idempotency ---

// CheckIdempotency returns true if this notification ID has already been processed.
func (c *Cache) CheckIdempotency(ctx context.Context, notificationID string) (bool, error) {
	key := fmt.Sprintf("idempotent:%s", notificationID)
	set, err := c.client.SetNX(ctx, key, "1", 24*time.Hour).Result()
	if err != nil {
		return false, err
	}
	return !set, nil // true means duplicate
}

func (c *Cache) Close() error {
	return c.client.Close()
}
