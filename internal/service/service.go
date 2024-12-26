package service

import (
	"Tasks/internal/interfaces"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	"Tasks/internal/model"
)

type Service struct {
	repo  interfaces.StorageRepository
	cache interfaces.CacheRepository
}

func NewService(repo interfaces.StorageRepository, repoCache interfaces.CacheRepository) *Service {
	return &Service{repo: repo, cache: repoCache}
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

func (s *Service) AddUser(ctx context.Context, userId int, taskID int) error {
	return s.repo.AddNewUserTask(ctx, userId, taskID)
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
	s.cache.UpdateTaskStatusInCache(ctx, taskID, newStatus)
	return s.repo.TaskUpdateStatus(ctx, newStatus, taskID)
}

func (s *Service) DeleteTask(ctx context.Context, taskID int) error {
	s.cache.DeleteTaskFromCache(ctx, taskID)
	return s.repo.DeleteTask(ctx, taskID)
}

func (s *Service) RemoveUserFromTask(ctx context.Context, userID int, taskID int) error {
	return s.repo.RemoveUserFromTask(ctx, userID, taskID)
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
