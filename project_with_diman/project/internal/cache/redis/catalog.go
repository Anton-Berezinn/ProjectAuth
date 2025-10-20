package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"time"
)

func NewClientCatalog(ctx context.Context) (*RedisCache, error) {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   2,
	})

	if err := client.Ping().Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	newClient := &RedisCache{client: client}
	return newClient, nil
}

func (r *RedisCache) InsertCatalog(ctx context.Context, data string, userId string) error {
	_, err := r.client.Set(userId, data, time.Hour*2).Result()
	if err != nil {
		return fmt.Errorf("failed to set data from redis server: %w", err)
	}
	return nil
}

func (r *RedisCache) UpdateCatalog(ctx context.Context, data string, userId string) error {
	_, err := r.client.Set(userId, data, time.Hour*2).Result()
	if err != nil {
		return fmt.Errorf("failed to set data from redis server: %w", err)
	}
	return nil
}

func (r *RedisCache) GetCatalogs(ctx context.Context, userId string) (string, error) {
	val, err := r.client.Get(userId).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found %w", ErrNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get data from redis server: %w", err)
	}
	return val, nil
}
