package main

import (
	repoCahe "Tasks/internal/repository/redis"
	"context"

	"Tasks/internal/app"
	"Tasks/internal/config"
	"Tasks/internal/http-server/handlers"
	"Tasks/internal/lib/logger"
	"Tasks/internal/lib/logger/sl"
	repo "Tasks/internal/repository/postgres"
	"Tasks/internal/service"
	"Tasks/internal/storage"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("start application")

	storages, err := storage.NewStorage(context.Background(), cfg)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
	}
	// todo: defer storage.Close()

	repoStorage := repo.NewStorage(storages.Postgres, log)
	repoCache := repoCahe.NewCache(storages.Redis, log)
	serv := service.NewService(repoStorage, repoCache)
	deps := &handlers.Dependencies{
		Service: serv,
		Log:     log,
	}

	h := handlers.NewHandler(deps)
	router := app.SetupRouter(h, log)
	// TODO: разобраться с авторизацией, первая попытка неудачна

	server := app.New(cfg, log, router)
	if err := server.Run(); err != nil {
		log.Error("server stopped with error", sl.Err(err))
	}
}
