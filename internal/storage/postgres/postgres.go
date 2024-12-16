package postgres

import (
	"TaskManagementSystemWithAnalytics/internal/config"
	"TaskManagementSystemWithAnalytics/internal/lib/logger/sl"
	"TaskManagementSystemWithAnalytics/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

// docker run --name=dbase -e POSTGRES_PASSWORD='1234' -e POSTGRES_DB='taskmanag' -p 5436:5432 -d --rm postgres
// psql -h localhost -U postgres -p 5436 -d  taskmanag

type Storage struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func New(ctx context.Context, log *slog.Logger, cfg *config.Config) (*Storage, error) {
	pool, err := pgxpool.New(ctx, cfg.DataBase.StorageURL)
	if err != nil {
		return nil, err
	}
	return &Storage{pool: pool, log: log}, nil
}

// методы create-new-task
// создание задачи
func (s *Storage) CreateNewTask(
	ctx context.Context,
	nameTask string,
	description string,
	deadline time.Time,
) (int, error) { // возвращаем taskID
	const op = "storage.postgres.CreateNewTask"
	log := s.log.With(slog.String("op", op))
	log.Info("create-new-task а new task")

	var taskID int

	createTask := "INSERT INTO tasks (name_task, description, deadline) " +
		"VALUES ($1, $2, $3) RETURNING task_id"

	err := s.pool.QueryRow(ctx, createTask, nameTask, description, deadline).Scan(&taskID)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return 0, fmt.Errorf("failed to create-new-task new task: %w", err)
	}
	log.Info("task created successfully", slog.Int("taskID", taskID))
	return taskID, nil
}

// методы read
// получение всех пользователей работающих над задачей
func (s *Storage) GetAllUsersWorkTask(ctx context.Context,
	taskID int) ([]model.User, error) {
	const op = "storage.postgres.GetAllUsersWorkTask"
	log := s.log.With(slog.String("op", op))
	log.Info("getting all the users working on the task")
	getUsers := "SELECT u.* FROM users u JOIN task_assignments ta ON u.user_id = ta.user_id WHERE ta.task_id = $1"

	rows, err := s.pool.Query(ctx, getUsers, taskID)
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
func (s *Storage) GetAllTasks(ctx context.Context, userID int) ([]model.Task, error) {
	const op = "storage.postgres.GetAllTasks"
	log := s.log.With(slog.String("op", op))
	log.Info("retrieving all tasks for user")

	getTasks := "SELECT t.* FROM tasks t JOIN task_assignments ta ON t.id = ta.task_id WHERE ta.user_id = $1;"

	rows, err := s.pool.Query(ctx, getTasks, userID)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Deadline, &task.Created_at, &task.Updated_at)
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
func (s *Storage) TaskShortDeadline(ctx context.Context) ([]model.Task, error) {
	const op = "storage.postgres.TaskShortDeadline"
	log := s.log.With(slog.String("op", op))
	log.Info("retrieving tasks with short deadlines")
	shortDeadline := "SELECT t.* FROM tasks t JOIN task_assignments ta ON t.id = ta.task_id WHERE ta.user_id = $1;"

	rows, err := s.pool.Query(ctx, shortDeadline)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()
	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Deadline, &task.Created_at, &task.Updated_at)
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
func (s *Storage) TaskUpdateStatus(ctx context.Context, newStatus string, taskID int) error {
	const op = "storage.postgres.TaskUpdateStatus"
	log := s.log.With(slog.String("op", op))
	log.Info("updating task status")
	query := "UPDATE tasks SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE task_id = $2"
	// Выполняем запрос к базе данных.
	_, err := s.pool.Exec(ctx, query, newStatus, taskID)
	if err != nil {
		log.Error("failed to update task status", sl.Err(err))
		return fmt.Errorf("failed to update task status: %w", err)
	}
	log.Info("task status updated successfully")
	return nil
}

// добавление пользователя к задачи
func (s *Storage) AddNewUserTask(ctx context.Context, userID int, taskID int) error { // статусы: to do, doing, done
	const op = "storage.postgres.AddNewUserTask"
	log := s.log.With(slog.String("op", op))
	log.Info("adding a user to a task")

	addUser := "INSERT INTO task_assignments (user_id, task_id) VALUES ($1, $2);"
	_, err := s.pool.Exec(ctx, addUser, userID, taskID)
	if err != nil {
		log.Error("failed to add user to task", sl.Err(err))
		return fmt.Errorf("failed to add user to task: %w", err)
	}
	log.Info("user successfully added to task")
	return nil
}

// методы delete
// удаление задачи
func (s *Storage) DeleteTask(ctx context.Context, taskID int) error {
	const op = "storage.postgres.DeleteTask"
	log := s.log.With(slog.String("op", op))

	log.Info("deleting a task")

	query := "DELETE FROM tasks WHERE task_id = $1"

	// Выполняем запрос к базе данных
	_, err := s.pool.Exec(ctx, query, taskID)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return fmt.Errorf("failed to delete task: %w", err)
	}

	log.Info("task deleted successfully")
	return nil
}

// Снятие пользователя с задачи
func (s *Storage) RemoveUserFromTask(ctx context.Context, userID int, taskID int) error {
	const op = "storage.postgres.RemoveUserFromTask"
	log := s.log.With(
		slog.String("op", op),
		slog.Int("userID", userID),
		slog.Int("taskID", taskID),
	)

	log.Info("removing user from task")

	query := "DELETE FROM task_assignments WHERE user_id = $1 AND task_id = $2"

	// Выполняем запрос к базе данных
	_, err := s.pool.Exec(ctx, query, userID, taskID)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return fmt.Errorf("failed to remove user from task: %w", err)
	}

	log.Info("user removed from task successfully")
	return nil
}

///
// TODO: разобраться с авторизацией, первая попытка неудачна
/*
func (s *Storage) SaveUser(login string, hashPass []byte) error {
	const op = "storage.postgres.Register"
	log := s.log.With(slog.String("op", op))
	log.Info("starting register new user")
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2);"

	_, err := s.pool.Exec(context.Background(), query, login, hashPass)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) SelectUser(login string) (model.User, error) {
	const op = "storage.postgres.Login"
	log := s.log.With(slog.String("op", op))
	log.Info("starting get user")
	var user model.User
	query := "SELECT user_id, username, password_hash, access_level FROM users WHERE username=$1;"

	err := s.pool.QueryRow(context.Background(), query, login).Scan(&user.ID, &user.Login, &user.HashPas, &user.Level)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("user not found", slog.String("login", login))
			return model.User{}, nil // Возвращаем nil, если пользователь не найден
		}
		log.Error("failed to query user", slog.String("login", login), slog.String("error", err.Error()))
		return model.User{}, err
	}

	log.Info("user found", slog.String("login", user.Login))
	return user, nil
}

*/
