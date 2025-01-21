package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"Tasks/internal/config"
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

func (s *Storage) Close(ctx context.Context) error {
	return s.Client.Close()
}
