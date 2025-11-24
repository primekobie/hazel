package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/primekobie/hazel/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateTask implements models.WorkspaceStore.
func (w *WorkspaceStore) CreateTask(ctx context.Context, task *models.Task) error {
	query := `INSERT INTO tasks(id, title, description, project_id, status, priority, due, created_at, last_modified)
	VALUES($1, $2, $3, $4, $5, $6, NULLIF($7,'0001-01-01 00:00:00'::TIMESTAMP), $8, $9);`

	_, err := w.conn.Exec(
		ctx,
		query,
		task.Id,
		task.Title,
		task.Description,
		task.Project.Id,
		task.Status,
		task.Priority,
		task.Due,
		task.CreatedAt,
		task.LastModified,
	)
	if err != nil {
		slog.Error("failed to insert task", "error", err.Error())
		return err
	}

	return nil
}

// DeleteTask implements models.WorkspaceStore.
func (w *WorkspaceStore) DeleteTask(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1;`

	_, err := w.conn.Exec(ctx, query, id)
	if err != nil {
		slog.Error("failed to delete task", "error", err.Error())
		return err
	}

	return nil
}

// GetTask implements models.WorkspaceStore.
func (w *WorkspaceStore) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	query := `SELECT
	t.id,
	t.title,
	t.description,
	t.status,
	t.priority,
	t.due,
	t.created_at,
	t.last_modified,
	p.id,
	p.name,
	p.description,
	p.created_at,
	p.last_modified
	FROM tasks AS t
	INNER JOIN projects AS p ON t.project_id = p.id
	WHERE t.id = $1;`

	task := &models.Task{Project: &models.Project{}}

	row := w.conn.QueryRow(ctx, query, id)
	err := row.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.Priority, &task.Due, &task.CreatedAt, &task.LastModified, &task.Project.Id, &task.Project.Name, &task.Project.Description, &task.Project.CreatedAt, &task.Project.LastModified)
	if err != nil {
		slog.Error("failed to scan task", "error", err.Error())
		return nil, err
	}

	return task, nil
}

// GetTasksForProject implements models.WorkspaceStore.
func (w *WorkspaceStore) GetTasksForProject(ctx context.Context, projectId uuid.UUID) ([]models.Task, error) {
	query := `SELECT
	id,
	title,
	description,
	status,
	priority,
	COALESCE(due,'0001-01-01 00:00:00'),
	created_at,
	last_modified
	FROM tasks
	WHERE project_id = $1;`

	tasks := []models.Task{}

	rows, err := w.conn.Query(ctx, query, projectId)
	if err != nil {
		slog.Error("failed to query tasks", "error", err.Error())
		return nil, err
	}

	for rows.Next() {
		var task models.Task

		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.Priority, &task.Due, &task.CreatedAt, &task.LastModified)
		if err != nil {
			slog.Error("failed to scan task", "error", err.Error())
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask implements models.WorkspaceStore.
func (w *WorkspaceStore) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `UPDATE tasks
	SET title = $1, description = $2, status = $3, priority = $4, due = $5, last_modified = $6
	WHERE id = $7;`

	_, err := w.conn.Exec(ctx, query, task.Title, task.Description, task.Status, task.Priority, task.Due, task.LastModified, task.Id)
	if err != nil {
		slog.Error("failed to scan task", "error", err.Error())
		return err
	}

	return nil
}

// AssignTask implements models.WorkspaceStore.
func (w *WorkspaceStore) AssignTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error {
	query := `INSERT INTO task_assignments(task_id, user_id)
	VALUES($1, $2)`

	_, err := w.conn.Exec(ctx, query, taskId, userId)
	if err != nil {
		slog.Error("failed to assign task to user", "error", err.Error())
		return err
	}

	return nil
}

// GetAssignedUsers implements models.WorkspaceStore.
func (w *WorkspaceStore) GetAssignedUsers(ctx context.Context, taskId uuid.UUID) ([]models.User, error) {
	query := `SELECT
	u.id,
	u.name,
	u.email,
	u.profile_photo,
	u.created_at,
	u.last_modified
	FROM task_assignments AS ta
	INNER JOIN users AS u ON ta.user_id = u.id
	WHERE ta.task_id = $1;`

	users := []models.User{}

	rows, err := w.conn.Query(ctx, query, taskId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}

		slog.Error("failed to fetch users", "error", err.Error())
		return nil, err
	}

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.ProfilePhoto, &user.CreatedAt, &user.LastModifed)
		if err != nil {
			slog.Error("failed to scan users", "error", err.Error())
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// UnassignTask implements models.WorkspaceStore.
func (w *WorkspaceStore) UnassignTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) error {
	query := `DELETE FROM task_assignments WHERE task_id = $1 AND user_id = $2;`

	_, err := w.conn.Exec(ctx, query, taskId, userId)
	if err != nil {
		slog.Error("failed to delete task assignment", "error", err.Error())
		return err
	}

	return nil
}
