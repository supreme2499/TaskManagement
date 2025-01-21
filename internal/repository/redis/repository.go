package repoCahe

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	redis2 "github.com/redis/go-redis/v9"

	"Tasks/internal/interfaces"
	"Tasks/internal/model"
	"Tasks/internal/storage/redis"
)

type Repo struct {
	redis *redis.Storage
	log   *slog.Logger
}

func NewCache(storage *redis.Storage, log *slog.Logger) interfaces.CacheRepository {
	return &Repo{redis: storage, log: log}
}

func (r *Repo) InsertingCache(ctx context.Context, task model.Task) error {
	const op = "repository.redis.InsertingCache"
	log := r.log.With(slog.String("op", op))
	log.Info("inserting the task")
	key := fmt.Sprintf("task:%d", task.ID)

	_, err := r.redis.Client.Pipelined(ctx, func(rdb redis2.Pipeliner) error {
		rdb.HSet(ctx, key, "NameTask", task.NameTask)
		rdb.HSet(ctx, key, "Description", task.Description)
		rdb.HSet(ctx, key, "Status", task.Status)
		rdb.HSet(ctx, key, "Deadline", task.Deadline.Format(time.RFC3339))
		rdb.HSet(ctx, key, "CreatedAt", task.CreatedAt.Format(time.RFC3339))
		rdb.HSet(ctx, key, "UpdatedAt", task.UpdatedAt.Format(time.RFC3339))
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to insert task into cache: %w", err)
	}
	return nil
}

func (r *Repo) GetTaskFromCache(ctx context.Context, taskID int) (model.Task, error) {
	const op = "repository.redis.GetTaskFromCache"
	log := r.log.With(slog.String("op", op))
	log.Info("getting task from cache")

	key := fmt.Sprintf("task:%d", taskID)

	fields, err := r.redis.Client.HGetAll(ctx, key).Result()
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to get task from cache: %w", err)
	}

	deadline, err := time.Parse(time.RFC3339, fields["Deadline"])
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to parse Deadline: %w", err)
	}
	createdAt, err := time.Parse(time.RFC3339, fields["CreatedAt"])
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to parse CreatedAt: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, fields["UpdatedAt"])
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to parse UpdatedAt: %w", err)
	}

	task := model.Task{
		ID:          taskID,
		NameTask:    fields["NameTask"],
		Description: fields["Description"],
		Status:      fields["Status"],
		Deadline:    deadline,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return task, nil
}

func (r *Repo) UpdateTaskStatusInCache(ctx context.Context, taskID int, status string) error {
	const op = "repository.redis.UpdateTaskStatusInCache"
	log := r.log.With(slog.String("op", op))
	log.Info("updating task status in cache")

	key := fmt.Sprintf("task:%d", taskID)

	_, err := r.redis.Client.HSet(ctx, key, "Status", status).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteTaskFromCache(ctx context.Context, taskID int) error {
	const op = "repository.redis.DeleteTaskFromCache"
	log := r.log.With(slog.String("op", op))
	log.Info("deleting task from cache")

	key := fmt.Sprintf("task:%d", taskID)

	_, err := r.redis.Client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
