package utils

import (
	"context"
	"log"
	"time"

	"github.com/anieswahdie1/ara-medika-api.git/internal/configs"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *configs.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Redis")
	return client, nil
}
