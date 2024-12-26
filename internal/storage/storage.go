package storage

import (
	"Tasks/internal/config"
	"Tasks/internal/storage/postgres"
	"Tasks/internal/storage/redis"
	"context"
)

type Storage struct {
	Postgres *postgres.Storage
	Redis    *redis.Storage
}

func NewStorage(ctx context.Context, cfg *config.Config) (*Storage, error) {
	// Инициализация PostgreSQL
	postgresStorage, err := postgres.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Инициализация Redis
	redisStorage, err := redis.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Storage{
		Postgres: postgresStorage,
		Redis:    redisStorage,
	}, nil
}
