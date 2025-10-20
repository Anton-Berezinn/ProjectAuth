package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"time"
)

var (
	ErrNotFound = errors.New("key not found")
)

type RedisCache struct {
	client *redis.Client
}

func NewClientUser(ctx context.Context) (*RedisCache, error) {
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
		DB:   1,
	})

	if err := client.Ping().Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	newClient := &RedisCache{client: client}
	return newClient, nil
}

func (r *RedisCache) GetData(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found %w", ErrNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("failed to get data from redis server: %w", err)
	}
	return val, nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found %w", ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("failed to delete data from redis server: %w", err)
	}
	return nil
}

func (r *RedisCache) UpdateData(ctx context.Context, key string, data int) error {
	val, err := r.client.Set(key, data, time.Hour*2).Result()
	if err != nil {
		return fmt.Errorf("failed to update data from redis server: %w", err)
	}
	fmt.Println(val)
	return nil
}

func (r *RedisCache) SetData(ctx context.Context, key string, val string) error {
	_, err := r.client.Set(key, val, time.Hour*2).Result()
	if err != nil {
		return fmt.Errorf("failed to set data from redis server: %w", err)
	}
	return nil

}
