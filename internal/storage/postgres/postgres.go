package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"Tasks/internal/config"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*Storage, error) {
	pdb, err := pgxpool.New(ctx, cfg.Postgres.StorageURL)
	if err != nil {
		return nil, err
	}
	return &Storage{Pool: pdb}, nil
}
