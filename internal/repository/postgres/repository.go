package repoStorage

import (
	"Tasks/internal/interfaces"
	"Tasks/internal/storage/postgres"
	"context"
	"fmt"
	"log/slog"
	"time"

	"Tasks/internal/lib/logger/sl"
	"Tasks/internal/model"
)

type Repo struct {
	postgres *postgres.Storage
	log      *slog.Logger
}

func NewStorage(storage *postgres.Storage, log *slog.Logger) interfaces.StorageRepository {
	return &Repo{postgres: storage, log: log}
}

// создание задачи
func (r *Repo) CreateNewTask(ctx context.Context, task model.Task) (int, error) { // возвращаем taskID
	const op = "storage.postgres.CreateNewTask"
	log := r.log.With(slog.String("op", op))
	log.Info("create-new-task а new task")

	query := "INSERT INTO tasks (name_task, description, deadline) " +
		"VALUES ($1, $2, $3) RETURNING task_id"
	err := r.postgres.Pool.QueryRow(ctx, query, task.NameTask, task.Description, task.Deadline).Scan(&task.ID)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return 0, fmt.Errorf("failed to create-new-task new task: %w", err)
	}
	log.Info("task created successfully", slog.Int("taskID", task.ID))
	return task.ID, nil
}

// получение всех пользователей работающих над задачей
func (r *Repo) GetAllUsersWorkTask(ctx context.Context, taskID int) ([]model.User, error) {
	const op = "storage.postgres.GetAllUsersWorkTask"
	log := r.log.With(slog.String("op", op))
	log.Info("getting all the users working on the task")
	getUsers := "SELECT u.* FROM users u JOIN task_assignments ta ON u.user_id = ta.user_id WHERE ta.task_id = $1"

	rows, err := r.postgres.Pool.Query(ctx, getUsers, taskID)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	var users []model.User
	var ignored string
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Login, &ignored, &user.Level)
		if err != nil {
			log.Error("failed to scan row", sl.Err(err))
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Error("row iteration error", sl.Err(err))
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	log.Info("successfully retrieved users", slog.Int("userCount", len(users)))
	return users, nil
}

// получение всех задач пользователя
func (r *Repo) GetAllTasks(ctx context.Context, userID int) ([]model.Task, error) {
	const op = "storage.postgres.GetAllTasks"
	log := r.log.With(slog.String("op", op))
	log.Info("retrieving all tasks for user")

	getTasks := "SELECT t.* FROM tasks t JOIN task_assignments ta ON t.task_id = ta.task_id WHERE ta.user_id = $1;"

	rows, err := r.postgres.Pool.Query(ctx, getTasks, userID)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.NameTask, &task.Description, &task.Status, &task.Deadline, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			log.Error("failed to scan row", sl.Err(err))
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tasks = append(tasks, task)

	}
	if err := rows.Err(); err != nil {
		log.Error("row iteration error", sl.Err(err))
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	log.Info("successfully retrieved tasks", slog.Int("taskCount", len(tasks)))
	return tasks, nil
}

// получение задач с приближающемся сроком
func (r *Repo) TaskShortDeadline(ctx context.Context, userID int) ([]model.Task, error) {
	const op = "storage.postgres.TaskShortDeadline"
	log := r.log.With(slog.String("op", op))
	log.Info("retrieving tasks with short deadlines")
	shortDeadline := `SELECT t.* 
                  FROM tasks t 
                  JOIN task_assignments ta ON t.task_id = ta.task_id 
                  WHERE ta.user_id = $1 
                  AND t.deadline BETWEEN $2 AND $3;`

	currentTime := time.Now()
	threeDaysLater := currentTime.Add(3 * 24 * time.Hour)

	rows, err := r.postgres.Pool.Query(ctx, shortDeadline, userID, currentTime, threeDaysLater)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()
	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.NameTask, &task.Description, &task.Status, &task.Deadline, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			log.Error("failed to scan row", sl.Err(err))
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tasks = append(tasks, task)
	}
	if rows.Err() != nil {
		log.Error("error encountered during row iteration", sl.Err(rows.Err()))
		return nil, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	log.Info("successfully retrieved tasks", slog.Int("taskCount", len(tasks)))
	return tasks, nil
}

// методы update
func (r *Repo) TaskUpdateStatus(ctx context.Context, newStatus string, taskID int) error {
	const op = "storage.postgres.TaskUpdateStatus"
	log := r.log.With(slog.String("op", op))
	log.Info("updating task status")
	query := "UPDATE tasks SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE task_id = $2"
	// Выполняем запрос к базе данных.
	_, err := r.postgres.Pool.Exec(ctx, query, newStatus, taskID)
	if err != nil {
		log.Error("failed to update task status", sl.Err(err))
		return fmt.Errorf("failed to update task status: %w", err)
	}
	log.Info("task status updated successfully")
	return nil
}

// добавление пользователя к задачи
func (r *Repo) AddNewUserTask(ctx context.Context, userID int, taskID int) error { // статусы: to do, doing, done
	const op = "storage.postgres.AddNewUserTask"
	log := r.log.With(slog.String("op", op))
	log.Info("adding a user to a task")

	addUser := "INSERT INTO task_assignments (user_id, task_id) VALUES ($1, $2);"
	_, err := r.postgres.Pool.Exec(ctx, addUser, userID, taskID)
	if err != nil {
		log.Error("failed to add user to task", sl.Err(err))
		return fmt.Errorf("failed to add user to task: %w", err)
	}
	log.Info("user successfully added to task")
	return nil
}

// методы delete
// удаление задачи
func (r *Repo) DeleteTask(ctx context.Context, taskID int) error {
	const op = "storage.postgres.DeleteTask"
	log := r.log.With(slog.String("op", op))

	log.Info("deleting a task")

	query := "DELETE FROM tasks WHERE task_id = $1"

	// Выполняем запрос к базе данных
	_, err := r.postgres.Pool.Exec(ctx, query, taskID)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return fmt.Errorf("failed to delete task: %w", err)
	}

	log.Info("task deleted successfully")
	return nil
}

// Снятие пользователя с задачи
func (r *Repo) RemoveUserFromTask(ctx context.Context, userID int, taskID int) error {
	const op = "storage.postgres.RemoveUserFromTask"
	log := r.log.With(
		slog.String("op", op),
		slog.Int("userID", userID),
		slog.Int("taskID", taskID),
	)

	log.Info("removing user from task")

	query := "DELETE FROM task_assignments WHERE user_id = $1 AND task_id = $2"

	// Выполняем запрос к базе данных
	_, err := r.postgres.Pool.Exec(ctx, query, userID, taskID)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return fmt.Errorf("failed to remove user from task: %w", err)
	}

	log.Info("user removed from task successfully")
	return nil
}

// Получение задачи по ID
func (r *Repo) TaskByID(ctx context.Context, taskID int) (model.Task, error) {
	const op = "storage.postgres.GetTaskByID"
	log := r.log.With(slog.String("op", op), slog.Int("taskID", taskID))

	log.Info("retrieving task by ID")

	query := "SELECT task_id, name_task, description, status, deadline, created_at, updated_at FROM tasks WHERE task_id = $1"

	var task model.Task
	err := r.postgres.Pool.QueryRow(ctx, query, taskID).Scan(
		&task.ID,
		&task.NameTask,
		&task.Description,
		&task.Status,
		&task.Deadline,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		if err.Error() == "sql: no rows in result set" {
			return model.Task{}, fmt.Errorf("task with ID %d not found", taskID)
		}
		return model.Task{}, fmt.Errorf("failed to retrieve task by ID: %w", err)
	}

	log.Info("task retrieved successfully", slog.Int("taskID", task.ID))
	return task, nil
}
