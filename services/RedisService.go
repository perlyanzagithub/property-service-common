package services

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"time"
)

type RedisService interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type redisService struct {
	client *redis.Client
}

// NewRedisService initializes the Redis service
func NewRedisService(client *redis.Client) RedisService {
	return &redisService{client: client}
}

func (r *redisService) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisService) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

func (r *redisService) Delete(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}
