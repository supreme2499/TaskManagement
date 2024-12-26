package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	resp "Tasks/internal/lib/api/response"
	"Tasks/internal/lib/logger/sl"
	"Tasks/internal/model"
	"Tasks/internal/service"
)

const invalid = "invalid request"

type Handler struct {
	service service.Service
	log     slog.Logger
}

type Dependencies struct {
	Service *service.Service
	Log     *slog.Logger
}

func NewHandler(deps *Dependencies) *Handler {
	return &Handler{
		service: *deps.Service,
		log:     *deps.Log,
	}
}

// Поступающие запросы
type RequestNewTask struct {
	TaskText    string    `json:"task_text" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Deadline    time.Time `json:"deadline" validate:"required"`
}

type RequestID struct {
	UserID int `json:"user_id" validate:"required"`
	TaskID int `json:"task_id" validate:"required"`
}

type RequestUserID struct {
	UserID int `json:"user_id" validate:"required"`
}

type RequestTaskID struct {
	TaskID int `json:"task_id" validate:"required"`
}

type RequestNewStatus struct {
	TaskID    int    `json:"task_id" validate:"required"`
	NewStatus string `json:"new_status" validate:"required"`
}

// Ответы
type ResponseNewTask struct {
	resp.Response
	TaskID int `json:"task_id"`
}

type Response struct {
	resp.Response
}

type ResponseTasks struct {
	Tasks []model.Task `json:"tasks"`
	resp.Response
}

type ResponseTask struct {
	Task model.Task `json:"task"`
	resp.Response
}

type ResponseUsers struct {
	Users []model.User `json:"users"`
	resp.Response
}

// Обработчики
func (h *Handler) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.CreateNewTask"
	log := h.log.With(slog.String("op", op))
	ctx := r.Context()
	req, err := decodeAndValidate[RequestNewTask](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}

	var task model.Task
	task.NameTask = req.TaskText
	task.Description = req.Description
	task.Deadline = req.Deadline
	taskID, err := h.service.CreateTask(ctx, task)
	if err != nil {
		errorHandler(log, "failed to create task", err, w, r)
		return
	}
	log.Info("task created successfully", slog.Int("task_id", taskID))
	render.JSON(w, r, ResponseNewTask{
		Response: resp.OK(),
		TaskID:   taskID,
	})
}

func (h *Handler) AddUserFromTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.AddUserFromTask"
	log := h.log.With(slog.String("op", op))
	ctx := r.Context()
	req, err := decodeAndValidate[RequestID](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}
	if err := h.service.AddUser(ctx, req.UserID, req.TaskID); err != nil {
		log.Error("failed to add user to task", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to add user to task"))
		return
	}
	log.Info("user added to task successfully", slog.Int("user_id", req.UserID), slog.Int("task_id", req.TaskID))
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}

func (h *Handler) AllTasks(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.AllTasks"
	log := h.log.With(
		slog.String("op", op))
	ctx := r.Context()
	req, err := decodeAndValidate[RequestUserID](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}
	tasks, err := h.service.AllTasks(ctx, req.UserID)
	if err != nil {
		log.Error("failed to retrieve tasks", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, resp.Error("failed to retrieve tasks"))
		return
	}
	log.Info("tasks retrieved successfully", slog.Int("user_id", req.UserID), slog.Int("task_count", len(tasks)))
	render.JSON(w, r, ResponseTasks{
		Tasks:    tasks,
		Response: resp.OK(),
	})
}

func (h *Handler) AllUsers(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.AllUsers"
	log := h.log.With(slog.String("op", op))
	ctx := r.Context()
	req, err := decodeAndValidate[RequestTaskID](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}
	users, err := h.service.AllUsersWorkTask(ctx, req.TaskID)
	if err != nil {
		errorHandler(log, "failed to retrieve users", err, w, r)
		return
	}
	log.Info("users retrieved successfully", slog.Int("task_id", req.TaskID), slog.Int("user_count", len(users)))
	render.JSON(w, r, ResponseUsers{
		Users:    users,
		Response: resp.OK(),
	})
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.DeleteTask"
	log := h.log.With(slog.String("op", op))
	ctx := r.Context()
	req, err := decodeAndValidate[RequestTaskID](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}
	if err := h.service.DeleteTask(ctx, req.TaskID); err != nil {
		errorHandler(log, "failed to delete task", err, w, r)
		return
	}
	log.Info("task deleted successfully", slog.Int("task_id", req.TaskID))
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}

func (h *Handler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.RemoveUser"
	log := h.log.With(slog.String("op", op))
	ctx := r.Context()
	req, err := decodeAndValidate[RequestID](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}
	if err := h.service.RemoveUserFromTask(ctx, req.UserID, req.TaskID); err != nil {
		errorHandler(log, "failed removing the user from the task", err, w, r)
		return
	}
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}

func (h *Handler) ShortDeadline(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.ShortDeadline"
	log := h.log.With(slog.String("op", op))
	ctx := r.Context()
	req, err := decodeAndValidate[RequestUserID](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}
	tasks, err := h.service.TaskShortDeadline(ctx, req.UserID)
	if err != nil {
		errorHandler(log, "failed gets tasks with a short deadline", err, w, r)
		return
	}
	render.JSON(w, r, ResponseTasks{
		Tasks:    tasks,
		Response: resp.OK(),
	})
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.UpdateStatus"
	log := h.log.With(slog.String("op", op))
	ctx := r.Context()
	req, err := decodeAndValidate[RequestNewStatus](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}
	if err := h.service.TaskUpdateStatus(ctx, req.NewStatus, req.TaskID); err != nil {
		errorHandler(log, "failed update task status", err, w, r)
		return
	}
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.GetTaskByID"
	log := h.log.With(slog.String("op", op))
	ctx := r.Context()
	req, err := decodeAndValidate[RequestTaskID](r, h.log)
	if err != nil {
		errorHandler(log, invalid, err, w, r)
		return
	}
	task, err := h.service.TaskByID(ctx, req.TaskID)
	if err != nil {
		errorHandler(log, "failed update task status", err, w, r)
		return
	}

	render.JSON(w, r, ResponseTask{
		Task:     task,
		Response: resp.OK(),
	})
}

// Вспомогательные функции
func decodeAndValidate[T any](r *http.Request, log slog.Logger) (*T, error) {
	var req T
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		return nil, fmt.Errorf("failed to decode request")
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", sl.Err(err))
		return nil, fmt.Errorf("validation error: %w", err)
	}
	return &req, nil
}

func errorHandler(log *slog.Logger, msg string, err error, w http.ResponseWriter, r *http.Request) {
	log.Error(msg, sl.Err(err))
	w.WriteHeader(http.StatusBadRequest)
	render.JSON(w, r, resp.Error(err.Error()))
	return
}
