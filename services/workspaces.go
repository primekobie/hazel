package services

import (
	"context"
	"time"

	"github.com/freekobie/hazel/models"
	"github.com/google/uuid"
)

type WorkspaceService struct {
	store models.WorkspaceStore
}

func NewWorkspaceService(store models.WorkspaceStore) *WorkspaceService {
	return &WorkspaceService{
		store: store,
	}
}

func (s *WorkspaceService) NewWorkspace(ctx context.Context, ws *models.Workspace) error {
	ws.Id = uuid.New()
	createdAt := time.Now().UTC()
	ws.CreatedAt = createdAt
	ws.LastModified = createdAt
	ws.User.Role = "admin"
	err := s.store.Create(ctx, ws)
	if err != nil {
		return err
	}

	return nil
}

func (s *WorkspaceService) GetWorkspace(ctx context.Context, id uuid.UUID) (*models.Workspace, error) {
	ws, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *WorkspaceService) GetUserWorkspaces(ctx context.Context, id uuid.UUID) ([]models.Workspace, error) {
	workspaces, err := s.store.GetAllForUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return workspaces, nil
}

func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, ws *models.Workspace) error {
	workspace, err := s.store.Get(ctx, ws.Id)
	if err != nil {
		return err
	}

	workspace.Name = ws.Name
	workspace.Description = ws.Description
	workspace.LastModified = time.Now()

	err = s.store.Update(ctx, workspace)
	if err != nil {
		return err
	}

	return nil
}
