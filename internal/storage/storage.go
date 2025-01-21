package storage

import (
	"context"
	"fmt"

	"Tasks/internal/config"
	"Tasks/internal/storage/postgres"
	"Tasks/internal/storage/redis"
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

func (s *Storage) Close(ctx context.Context) error {
	var closeErrs []error

	// Закрытие PostgreSQL
	if err := s.Postgres.Close(); err != nil {
		closeErrs = append(closeErrs, fmt.Errorf("failed to close Postgres: %w", err))
	}

	// Закрытие Redis
	if err := s.Redis.Close(ctx); err != nil {
		closeErrs = append(closeErrs, fmt.Errorf("failed to close Redis: %w", err))
	}

	if len(closeErrs) > 0 {
		return fmt.Errorf("storage close errors: %v", closeErrs)
	}

	return nil
}
