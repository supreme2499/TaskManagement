package main

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"Tasks/internal/config"
)

// docker run --name=dbase -e POSTGRES_PASSWORD='1234' -e POSTGRES_DB='taskmanag' -p 5436:5432 -d --rm postgres
func main() {
	cfg := config.MustLoad()
	var migrationsPath = cfg.Postgres.MigrationsPath
	if migrationsPath == "" {
		panic("MIGRATIONS_PATH is required")
	}
	m, err := migrate.New("file://"+migrationsPath, cfg.Postgres.StorageURL)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic(err)
	}
	fmt.Println("migrations applied successfully")
}
