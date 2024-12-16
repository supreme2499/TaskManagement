package create_new_task

import (
	resp "TaskManagementSystemWithAnalytics/internal/lib/api/response"
	"TaskManagementSystemWithAnalytics/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"time"
)

import (
	"context"
)

type Request struct {
	TaskText    string    `json:"task_text" validate:"required,TaskText"`
	Description string    `json:"description" validate:"required,Description"`
	Deadline    time.Time `json:"deadline" validate:"required,Deadline"`
}

type Response struct {
	resp.Response
	TaskID int `json:"task_id"`
}

type TaskSaver interface {
	CreateNewTask(ctx context.Context, nameTask string, description string, deadline time.Time) (int, error)
}

func New(ctx context.Context, log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "management.create-new-task"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		taskID, err := taskSaver.CreateNewTask(ctx, req.TaskText, req.Description, req.Deadline)
		if err != nil {
			log.Error("failed to create-new-task new task", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to create-new-task new task"))
			return
		}
		render.JSON(w, r, Response{
			Response: resp.OK(),
			TaskID:   taskID,
		})
	}
}
