package service

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore interface {
	Exists(ctx context.Context, key string) (int64, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	ZIncrBy(ctx context.Context, key string, increment float64, member string) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	SCard(ctx context.Context, key string) (int64, error)
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
}

type RedisStoreAdapter struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStoreAdapter {
	return &RedisStoreAdapter{client: client}
}

func (r *RedisStoreAdapter) Exists(ctx context.Context, key string) (int64, error) {
	return r.client.Exists(ctx, key).Result()
}

func (r *RedisStoreAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisStoreAdapter) ZIncrBy(ctx context.Context, key string, increment float64, member string) error {
	return r.client.ZIncrBy(ctx, key, increment, member).Err()
}

func (r *RedisStoreAdapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *RedisStoreAdapter) SCard(ctx context.Context, key string) (int64, error) {
	return r.client.SCard(ctx, key).Result()
}

func (r *RedisStoreAdapter) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

func (r *RedisStoreAdapter) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

func (r *RedisStoreAdapter) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}
