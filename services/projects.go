package services

import (
	"context"
	"time"

	"github.com/freekobie/hazel/models"
	"github.com/google/uuid"
)

func (s *WorkspaceService) CreateProject(ctx context.Context, project *models.Project) error {
	project.Id = uuid.New()

	lastModified := time.Now()
	project.CreatedAt = lastModified
	project.LastModified = lastModified
	project.Status = "active"

	err := s.store.CreateProject(ctx, project)
	if err != nil {
		return err
	}

	return nil
}

func (s *WorkspaceService) GetProject(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.store.GetProject(ctx, id)
}

func (s *WorkspaceService) UpdateProject(ctx context.Context, data map[string]any) (*models.Project, error) {

	id, _ := data["id"]
	project, err := s.store.GetProject(ctx, id.(uuid.UUID))
	if err != nil {
		return nil, err
	}

	name, ok := data["name"]
	if ok {
		project.Name = name.(string)
	}

	description, ok := data["description"]
	if ok {
		project.Description = description.(string)
	}

	startDate, ok := data["startDate"]
	if ok {
		newTime, err := time.Parse(models.DateLayout, startDate.(string))
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		project.StartDate.Time = newTime
	}

	endDate, ok := data["endDate"]
	if ok {
		newTime, err := time.Parse(models.DateLayout, endDate.(string))
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		project.EndDate.Time = newTime
	}

	project.LastModified = time.Now()

	err = s.store.UpdateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *WorkspaceService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return s.store.DeleteProject(ctx, id)
}
