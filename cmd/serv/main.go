package main

import (
	"TaskManagementSystemWithAnalytics/internal/config"
	create "TaskManagementSystemWithAnalytics/internal/http-server/handlers/manager-analitic/create-new-task"
	delete_task "TaskManagementSystemWithAnalytics/internal/http-server/handlers/manager-analitic/delete/delete-task"
	add_user "TaskManagementSystemWithAnalytics/internal/http-server/handlers/manager-analitic/update/add-user"
	update_status "TaskManagementSystemWithAnalytics/internal/http-server/handlers/manager-analitic/update/update-status"
	mwLogger "TaskManagementSystemWithAnalytics/internal/http-server/middleware/logger"
	"TaskManagementSystemWithAnalytics/internal/lib/logger/handler/slogpretty"
	"TaskManagementSystemWithAnalytics/internal/lib/logger/sl"
	"TaskManagementSystemWithAnalytics/internal/storage/postgres"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// $env:CONFIG_PATH="./config/local.yaml"; $env:SECRET="supreme2499"; go run ./cmd/serv/main.go
// $env:CONFIG_PATH="./config/local.yaml"; go run ./cmd/migrator/main.go

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("start application")

	storage, err := postgres.New(context.Background(), log, cfg)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
	}
	_ = storage

	router := chi.NewRouter()
	ctx := context.Background()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: разобраться с авторизацией, первая попытка неудачна
	/*
		router.Put("/register", auth.Register(log, storage))
		router.Post("/login", auth.Login(log, storage))
	*/
	//router.Route("/protected", func(r chi.Router) {
	//	// TODO: исправить авторизацию, пока что овер плохо, но это так, для тестов
	//	r.Use(middleware.BasicAuth("url-shortener", map[string]string{
	//		"surp": "1234",
	//	}))
	//
	//	r.Post("/", create-new-task.New(ctx, log, storage))
	//})

	router.Post("/task", create.New(ctx, log, storage))
	router.Put("/upu", add_user.New(ctx, log, storage))
	router.Put("/ups", update_status.New(ctx, log, storage))
	router.Delete("/task", delete_task.New(ctx, log, storage))

	log.Info("starting server", slog.String("address", cfg.Http.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Http.Address,
		Handler:      router,
		ReadTimeout:  cfg.Http.Timeout,
		WriteTimeout: cfg.Http.Timeout,
		IdleTimeout:  cfg.Http.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	// TODO: close storage

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
