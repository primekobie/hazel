package services

import (
	"context"
	"strings"
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
	ws.User.Role = "owner"

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

func (s *WorkspaceService) UpdateWorkspace(ctx context.Context, wsData map[string]string) (*models.Workspace, error) {
	id, _ := wsData["id"]
	workspace, err := s.store.Get(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, err
	}

	name, ok := wsData["name"]
	if ok {
		workspace.Name = name
	}

	description, ok := wsData["description"]
	if ok {
		workspace.Description = description
	}

	workspace.LastModified = time.Now()

	err = s.store.Update(ctx, workspace)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (s *WorkspaceService) DeleteWorkspace(ctx context.Context, id uuid.UUID) error {
	return s.store.Delete(ctx, id)
}

func (s *WorkspaceService) AddWorkspaceMember(ctx context.Context, workspaceId, userId uuid.UUID, role string) error {
	err := s.store.AddMembership(ctx, workspaceId, userId, role)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return ErrDuplicateEntry
		}

		return ErrFailedOperation
	}

	return nil
}

func (s *WorkspaceService) GetWorkspaceMembers(ctx context.Context, id uuid.UUID) (any, error) {
	return s.store.GetWorkspaceMembers(ctx, id)
}

func (s *WorkspaceService) DeleteWorkspaceMember(ctx context.Context, workspaceId, userId uuid.UUID) error {
	return s.store.DeleteMembership(ctx, workspaceId, userId)
}
