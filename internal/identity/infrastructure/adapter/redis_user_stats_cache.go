package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmlc643/twitbet-backend/internal/identity/domain/entity"
	"github.com/redis/go-redis/v9"
)

type RedisUserStatsCache struct {
	rdb *redis.Client
}

func NewRedisUserStatsCache(rdb *redis.Client) *RedisUserStatsCache {
	return &RedisUserStatsCache{
		rdb: rdb,
	}
}

func (c *RedisUserStatsCache) GetStats(ctx context.Context, userID string) (*entity.UserStats, error) {
	cacheKey := fmt.Sprintf("user_stats:%s", userID)

	cachedStats, err := c.rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var stats entity.UserStats
	if err := json.Unmarshal([]byte(cachedStats), &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (c *RedisUserStatsCache) SetStats(ctx context.Context, userID string, stats *entity.UserStats, expiration time.Duration) error {
	cacheKey := fmt.Sprintf("user_stats:%s", userID)

	statsBytes, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, cacheKey, statsBytes, expiration).Err()
}
