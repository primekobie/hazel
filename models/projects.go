package models

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Date struct {
	time.Time
}

const DateLayout = "2006-01-02"

func (d *Date) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", d.Format(DateLayout))
	return []byte(formatted), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == "" {
		d.Time = time.Time{}
		return nil
	}
	if len(str) > 0 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	parsed, err := time.Parse(DateLayout, str)
	if err != nil {
		slog.Error("failed to parse time", "error", err.Error())
		return err
	}
	d.Time = parsed
	return nil
}

type Project struct {
	Id           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Workspace    *Workspace `json:"workspace,omitempty"`
	StartDate    Date       `json:"startDate,omitzero"`
	EndDate      Date       `json:"endDate,omitzero"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	LastModified time.Time  `json:"lastModified"`
}

type ProjectStore interface {
	CreateProject(ctx context.Context, project *Project) error
	UpdateProject(ctx context.Context, project *Project) error
	GetProject(ctx context.Context, id uuid.UUID) (*Project, error)
	GetWorkspaceProjects(ctx context.Context, workspaceId uuid.UUID) ([]Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID) error
}
