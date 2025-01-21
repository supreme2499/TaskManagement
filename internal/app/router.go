package app

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"Tasks/internal/http-server/handlers"
	mwLogger "Tasks/internal/http-server/middleware/logger"
)

func SetupRouter(h *handlers.Handler, log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/task", h.CreateNewTask)
	router.Post("/adduser", h.AddUserFromTask)
	router.Get("/users", h.AllUsers)
	router.Get("/tasks", h.AllTasks)
	router.Get("/shortdeadline", h.ShortDeadline)
	router.Get("/taskbyid", h.GetTaskByID)
	router.Put("/status", h.UpdateStatus)
	router.Delete("/task", h.DeleteTask)
	router.Delete("/user", h.RemoveUser)

	return router
}
