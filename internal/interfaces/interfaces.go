package interfaces

import (
	"context"

	"Tasks/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@v2.51.1 --name=StorageRepository --output=../service/mocks
type StorageRepository interface {
	CreateNewTask(ctx context.Context, task model.Task) (int, error)
	GetAllUsersWorkTask(ctx context.Context, taskID int) ([]model.User, error)
	GetAllTasks(ctx context.Context, userID int) ([]model.Task, error)
	TaskShortDeadline(ctx context.Context, userID int) ([]model.Task, error)
	TaskUpdateStatus(ctx context.Context, newStatus string, taskID int) error
	AddNewUserTask(ctx context.Context, userID int, taskID int) error
	DeleteTask(ctx context.Context, taskID int) error
	RemoveUserFromTask(ctx context.Context, userID int, taskID int) error
	TaskByID(ctx context.Context, taskID int) (model.Task, error)
	UserByID(ctx context.Context, taskID int) ([]int, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.51.1 --name=CacheRepository --output=../service/mocks
type CacheRepository interface {
	InsertingCache(ctx context.Context, task model.Task) error
	GetTaskFromCache(ctx context.Context, taskID int) (model.Task, error)
	UpdateTaskStatusInCache(ctx context.Context, taskID int, status string) error
	DeleteTaskFromCache(ctx context.Context, taskID int) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.51.1 --name=Broker --output=../service/mocks
type Broker interface {
	Produce(message []byte, topic string) error
}
