package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/freekobie/hazel/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (w *WorkspaceStore) CreateProject(ctx context.Context, project *models.Project) error {
	query := `INSERT INTO projects(id, name, description, workspace_id, start_date, end_date, created_at, last_modified)
	VALUES($1, $2, $3, $4, NULLIF($5,'0001-01-01'::DATE),NULLIF($6,'0001-01-01'::DATE), $7, $8);`

	_, err := w.conn.Exec(
		ctx,
		query,
		project.Id,
		project.Name,
		project.Description,
		project.Workspace.Id,
		project.StartDate.Format(models.DateLayout),
		project.EndDate.Format(models.DateLayout),
		project.CreatedAt,
		project.LastModified,
	)
	if err != nil {
		slog.Error("failed to insert project", "error", err.Error())
		return err
	}

	return nil
}

func (w *WorkspaceStore) UpdateProject(ctx context.Context, project *models.Project) error {
	query := `UPDATE projects SET name = $1, description = $2, start_date = NULLIF($3,'0001-01-01'::DATE), end_date = NULLIF($4,'0001-01-01'::DATE), last_modified = $5 WHERE id = $6;`
	_, err := w.conn.Exec(
		ctx,
		query,
		project.Name,
		project.Description,
		project.StartDate.Format(models.DateLayout),
		project.EndDate.Format(models.DateLayout),
		project.LastModified,
		project.Id,
	)
	if err != nil {
		slog.Error("failed to update project", "error", err.Error())
		return err
	}

	return nil
}

func (w *WorkspaceStore) GetProject(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	query := `SELECT
	p.id,
	p.name,
	p.description,
	COALESCE(p.start_date,'0001-01-01'),
	COALESCE(p.end_date,'0001-01-01'),
	p.created_at,
	p.last_modified,
	w.id,
	w.name,
	w.description,
	w.created_at,
	w.last_modified
	FROM projects AS p
	INNER JOIN workspaces AS w ON p.workspace_id = w.id
	WHERE p.id = $1;`

	row := w.conn.QueryRow(ctx, query, id)
	project := &models.Project{Workspace: &models.Workspace{}}

	err := row.Scan(&project.Id, &project.Name, &project.Description, &project.StartDate.Time, &project.EndDate.Time, &project.CreatedAt, &project.LastModified, &project.Workspace.Id, &project.Workspace.Name, &project.Workspace.Description, &project.Workspace.CreatedAt, &project.Workspace.LastModified)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		slog.Error("failed to read project", "error", err.Error())
		return nil, err
	}

	return project, nil
}

func (w *WorkspaceStore) GetWorkspaceProjects(ctx context.Context, workspaceId uuid.UUID) ([]models.Project, error) {
	return nil, nil
}

func (w *WorkspaceStore) DeleteProject(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1;`

	_, err := w.conn.Exec(ctx, query, id)
	if err != nil {
		slog.Error("failed to delete project", "error", err.Error())
		return err
	}
	return nil
}
