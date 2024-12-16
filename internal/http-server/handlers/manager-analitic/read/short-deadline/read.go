package short_deadline

import (
	resp "TaskManagementSystemWithAnalytics/internal/lib/api/response"
	"TaskManagementSystemWithAnalytics/internal/lib/logger/sl"
	"TaskManagementSystemWithAnalytics/internal/model"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	Tasks []model.Task `json:"users"`
	resp.Response
}

type TaskSaver interface {
	TaskShortDeadline(ctx context.Context) ([]model.Task, error)
}

func New(ctx context.Context, log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.manager-analitic.read.short-deadline"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		tasks, err := taskSaver.TaskShortDeadline(ctx)
		if err != nil {
			log.Error("failed to create-new-task new task", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to create-new-task new task"))
			return
		}
		render.JSON(w, r, Response{
			Tasks:    tasks,
			Response: resp.OK(),
		})
	}
}
