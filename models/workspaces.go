package models

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("entity not found")
)

// Workspace represents a top-level organizational unit or collaboration space.
// Projects and Users belong to a Workspace.
type Workspace struct {
	Id           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	User         *User     `json:"user,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`
}

type WorkspaceStore interface {
	Create(ctx context.Context, workspace *Workspace) error
	Update(ctx context.Context, workspace *Workspace) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*Workspace, error)
	GetAllForUser(ctx context.Context, userId uuid.UUID) ([]Workspace, error)
	GetWorkspaceMembers(ctx context.Context, workspaceId uuid.UUID) ([]User, error)
	AddMembership(ctx context.Context, workspaceId, userId uuid.UUID, role string) error
	DeleteMembership(ctx context.Context, workspaceId, userId uuid.UUID) error
	ProjectStore
	TaskStore
}
