package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"Tasks/internal/config"
)

type Server struct {
	cfg    *config.Config
	log    *slog.Logger
	server *http.Server
}

func New(cfg *config.Config, log *slog.Logger, handler http.Handler) *Server {
	return &Server{
		cfg: cfg,
		log: log,
		server: &http.Server{
			Addr:         cfg.HTTP.Address,
			Handler:      handler,
			ReadTimeout:  cfg.HTTP.Timeout,
			WriteTimeout: cfg.HTTP.Timeout,
			IdleTimeout:  cfg.HTTP.IdleTimeout,
		},
	}
}

func (s *Server) Run() error {
	s.log.Info("starting server", slog.String("address", s.cfg.HTTP.Address))

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("server error", slog.String("error", err.Error()))
		}
	}()

	s.log.Info("server started")

	// Ожидание сигнала завершения
	<-doneCh
	s.log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.HTTP.WithTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("failed to shutdown server", slog.String("error", err.Error()))
		return err
	}

	s.log.Info("server stopped gracefully")
	return nil
}
