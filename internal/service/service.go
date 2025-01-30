package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"Tasks/internal/interfaces"
	"Tasks/internal/lib/logger/sl"
	"Tasks/internal/model"
)

type Service struct {
	log      *slog.Logger
	repo     interfaces.StorageRepository
	cache    interfaces.CacheRepository
	producer interfaces.Broker
}

func NewService(log *slog.Logger,
	repo interfaces.StorageRepository,
	repoCache interfaces.CacheRepository,
	producer interfaces.Broker) *Service {
	return &Service{log: log, repo: repo, cache: repoCache, producer: producer}
}

func (s *Service) CreateTask(ctx context.Context, task model.Task) (int, error) {
	currentTime := time.Now()
	if task.Deadline.Before(currentTime) {
		return 0, fmt.Errorf("task deadline is too far in the past")
	}

	taskID, err := s.repo.CreateNewTask(ctx, task)
	if err != nil {
		return 0, err
	}
	err = s.cache.InsertingCache(ctx, task)
	if err != nil {
		return taskID, fmt.Errorf("cache insertion failed: %w", err)
	}
	return taskID, nil
}

func (s *Service) AddUser(ctx context.Context, userID int, taskID int) error {
	err := s.repo.AddNewUserTask(ctx, userID, taskID)
	if err != nil {
		return err
	}

	msg := model.NotificationMessage{
		Event:     "add_user",
		Timestamp: time.Now().UTC(),
		TaskID:    taskID,
		UserID:    userID,
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = s.producer.Produce(msgJSON, "notification")
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	return nil
}

func (s *Service) AllUsersWorkTask(ctx context.Context, taskID int) ([]model.User, error) {
	return s.repo.GetAllUsersWorkTask(ctx, taskID)
}

func (s *Service) AllTasks(ctx context.Context, userID int) ([]model.Task, error) {
	return s.repo.GetAllTasks(ctx, userID)
}

func (s *Service) TaskShortDeadline(ctx context.Context, userID int) ([]model.Task, error) {
	return s.repo.TaskShortDeadline(ctx, userID)
}

func (s *Service) TaskUpdateStatus(ctx context.Context, newStatus string, taskID int) error {
	const op = "service.TaskUpdateStatus"
	log := s.log.With(slog.String("op", op))
	err := s.cache.UpdateTaskStatusInCache(ctx, taskID, newStatus)
	if err != nil {
		return err
	}
	err = s.repo.TaskUpdateStatus(ctx, newStatus, taskID)
	if err != nil {
		return err
	}

	users, err := s.repo.UserByID(ctx, taskID)
	if err != nil {
		return err
	}

	for _, user := range users {
		msg := model.NotificationMessage{
			Event:        "change_status",
			Timestamp:    time.Now().UTC(),
			TaskID:       taskID,
			UserID:       user,
			ChangeStatus: newStatus,
		}

		msgJSON, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}
		err = s.producer.Produce(msgJSON, "notification")
		if err != nil {
			log.Error("failed to produce message", sl.Err(err))
		}
	}
	return nil
}

func (s *Service) DeleteTask(ctx context.Context, taskID int) error {
	const op = "service.DeleteTask"
	log := s.log.With(slog.String("op", op))

	err := s.cache.DeleteTaskFromCache(ctx, taskID)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteTask(ctx, taskID); err != nil {
		return err
	}

	users, err := s.repo.UserByID(ctx, taskID)
	if err != nil {
		return err
	}

	for _, user := range users {
		msg := model.NotificationMessage{
			Event:     "delete_task",
			Timestamp: time.Now().UTC(),
			TaskID:    taskID,
			UserID:    user,
		}

		msgJSON, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}
		err = s.producer.Produce(msgJSON, "notification")
		if err != nil {
			log.Error("failed to produce message", sl.Err(err))
		}
	}
	return nil
}

func (s *Service) RemoveUserFromTask(ctx context.Context, userID int, taskID int) error {
	log := s.log.With(slog.String("op", "service.RemoveUserFromTask"))
	err := s.repo.RemoveUserFromTask(ctx, userID, taskID)
	if err != nil {
		return err
	}
	log.Info("successful removal of the user from the database")

	msg := model.NotificationMessage{
		Event:     "remove_user",
		Timestamp: time.Now().UTC(),
		TaskID:    taskID,
		UserID:    userID,
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = s.producer.Produce(msgJSON, "notifications")
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	log.Info("successfully sending a message to kafka")
	return nil
}

func (s *Service) TaskByID(ctx context.Context, taskID int) (model.Task, error) {
	task, err := s.cache.GetTaskFromCache(ctx, taskID)
	if err == nil {
		return task, nil
	}
	if errors.Is(err, redis.Nil) {
		task, err = s.repo.TaskByID(ctx, taskID)
		if err != nil {
			return model.Task{}, err
		}
		err = s.cache.InsertingCache(ctx, task)
		if err != nil {
			return task, fmt.Errorf("cache insertion failed: %w", err)
		}
	}
	return task, nil
}
