package all_users

import (
	resp "TaskManagementSystemWithAnalytics/internal/lib/api/response"
	"TaskManagementSystemWithAnalytics/internal/lib/logger/sl"
	"TaskManagementSystemWithAnalytics/internal/model"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	TaskID int `json:"task_id" validate:"required, task_id"`
}

type Response struct {
	Users []model.User `json:"users"`
	resp.Response
}

type TaskSaver interface {
	GetAllUsersWorkTask(ctx context.Context, taskID int) ([]model.User, error)
}

func New(ctx context.Context, log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.manager-analitic.read.all-users"
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
		users, err := taskSaver.GetAllUsersWorkTask(ctx, req.TaskID)
		if err != nil {
			log.Error("failed to create-new-task new task", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to create-new-task new task"))
			return
		}
		render.JSON(w, r, Response{
			Users:    users,
			Response: resp.OK(),
		})
	}
}
