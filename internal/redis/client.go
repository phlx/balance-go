package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type redisClient struct {
	internalClient *redis.Client
}

func (rc redisClient) Get(ctx context.Context, key string) (string, error) {
	return rc.internalClient.Get(ctx, key).Result()
}

func (rc redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.internalClient.Set(ctx, key, value, expiration).Err()
}

func New(ctx context.Context, addr string) (Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to ping Redis")
	}

	r := redisClient{internalClient: client}

	return r, nil
}
