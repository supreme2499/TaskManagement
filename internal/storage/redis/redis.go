package redis

import (
	"Tasks/internal/config"
	"context"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Client *redis.Client
}

func New(ctx context.Context, cfg *config.Config) (*Storage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
		DB:   cfg.Redis.Database,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return &Storage{Client: rdb}, nil
}
