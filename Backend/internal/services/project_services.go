package service

import (
	"context"
	"errors"
	"taskflow/internal/models"
	"taskflow/internal/repository"
	"time"

	"github.com/google/uuid"
)

type ProjectService struct {
	projectRepo *repository.ProjectRepository
}

func NewProjectService(projectRepo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, name, description string, ownerID uuid.UUID) error {

	if name == "" {
		return errors.New("project name is required")
	}

	project := &models.Project{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
	}

	return s.projectRepo.CreateProject(ctx, project)
}

func (s *ProjectService) GetProjects(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]models.Project, error) {
	return s.projectRepo.GetProjectsByUser(ctx, userID, limit, offset)
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, errors.New("project not found")
	}

	return project, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id uuid.UUID, name, description string, userID uuid.UUID) error {

	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	// 🔐 authorization check
	if project.OwnerID != userID {
		return errors.New("forbidden")
	}

	project.Name = name
	project.Description = description

	return s.projectRepo.UpdateProject(ctx, project)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {

	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	// 🔐 authorization check
	if project.OwnerID != userID {
		return errors.New("forbidden")
	}

	return s.projectRepo.DeleteProject(ctx, id)
}
