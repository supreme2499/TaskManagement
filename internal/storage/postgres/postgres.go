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

// docker run --name=dbase -e POSTGRES_PASSWORD='1234' -e POSTGRES_DB='taskmanag' -p 5436:5432 -d --rm postgres
// psql -h localhost -U postgres -p 5436 -d  taskmanag

// TODO: разобраться с авторизацией, первая попытка неудачна
/*
func (s *Storage) SaveUser(login string, hashPass []byte) error {
	const op = "storage.postgres.Register"
	log := s.log.With(slog.String("op", op))
	log.Info("starting register new user")
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2);"

	_, err := s.pool.Exec(context.Background(), query, login, hashPass)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) SelectUser(login string) (model.User, error) {
	const op = "storage.postgres.Login"
	log := s.log.With(slog.String("op", op))
	log.Info("starting get user")
	var user model.User
	query := "SELECT user_id, username, password_hash, access_level FROM users WHERE username=$1;"

	err := s.pool.QueryRow(context.Background(), query, login).Scan(&user.ID, &user.Login, &user.HashPas, &user.Level)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("user not found", slog.String("login", login))
			return model.User{}, nil // Возвращаем nil, если пользователь не найден
		}
		log.Error("failed to query user", slog.String("login", login), slog.String("error", err.Error()))
		return model.User{}, err
	}

	log.Info("user found", slog.String("login", user.Login))
	return user, nil
}

*/
