package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("error parseando la URL de Redis: %w", err)
	}

	rdb := redis.NewClient(opt)

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error al hacer Ping a Redis: %w", err)
	}
	
	log.Println("Conexión exitosa a Redis")
	return rdb, nil
}