package adapter

import (
	"context"
	"errors"
	"time"

	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/redis/go-redis/v9"
)

type RedisOTPRepository struct {
	client *redis.Client
}

func NewRedisOTPRepository(client *redis.Client) port.OTPRepository {
	return &RedisOTPRepository{
		client: client,
	}
}

func (r *RedisOTPRepository) SaveOTP(ctx context.Context, key, otpCode string, ttl time.Duration) error {
	return r.client.Set(ctx, key, otpCode, ttl).Err()
}

func (r *RedisOTPRepository) GetOTP(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return val, err
}

func (r *RedisOTPRepository) DeleteOTP(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
