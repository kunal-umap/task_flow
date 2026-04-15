package service

import (
	"context"
	"errors"
	"taskflow/internal/models"
	"taskflow/internal/repository"
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo    *repository.TaskRepository
	projectRepo *repository.ProjectRepository
}

func NewTaskService(taskRepo *repository.TaskRepository, projectRepo *repository.ProjectRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	title string,
	description string,
	projectID uuid.UUID,
	assigneeID *uuid.UUID,
	priority models.TaskPriority,
	dueDate *time.Time,
) error {

	if title == "" {
		return errors.New("title is required")
	}

	// check project exists
	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	// validate priority
	if priority != models.PriorityLow &&
		priority != models.PriorityMedium &&
		priority != models.PriorityHigh {
		return errors.New("invalid priority")
	}

	task := &models.Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      models.StatusTodo,
		Priority:    priority,
		ProjectID:   projectID,
		AssigneeID:  assigneeID,
		DueDate:     dueDate,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return s.taskRepo.CreateTask(ctx, task)
}

func (s *TaskService) GetTasks(
	ctx context.Context,
	projectID uuid.UUID,
	status *models.TaskStatus,
	assigneeID *uuid.UUID,
	limit int,
	offset int,
) ([]models.Task, error) {

	project, err := s.projectRepo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	return s.taskRepo.GetTasksByProject(ctx, projectID, status, assigneeID, limit, offset)
}

func (s *TaskService) GetTaskByID(
	ctx context.Context,
	taskID uuid.UUID,
	userID uuid.UUID,
) (*models.Task, error) {

	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	project, err := s.projectRepo.GetProjectByID(ctx, task.ProjectID)
	if err != nil {
		return nil, err
	}

	// 🔐 authorization: user must be owner OR assignee
	if project.OwnerID != userID &&
		(task.AssigneeID == nil || *task.AssigneeID != userID) {
		return nil, errors.New("forbidden")
	}

	return task, nil
}

func (s *TaskService) UpdateTask(
	ctx context.Context,
	taskID uuid.UUID,
	title string,
	description string,
	status models.TaskStatus,
	priority models.TaskPriority,
	assigneeID *uuid.UUID,
	dueDate *time.Time,
	userID uuid.UUID,
) error {

	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	project, err := s.projectRepo.GetProjectByID(ctx, task.ProjectID)
	if err != nil {
		return err
	}

	// 🔐 authorization (only project owner)
	if project.OwnerID != userID {
		return errors.New("forbidden")
	}

	// validate status
	if status != models.StatusTodo &&
		status != models.StatusInProgress &&
		status != models.StatusDone {
		return errors.New("invalid status")
	}

	// validate priority
	if priority != models.PriorityLow &&
		priority != models.PriorityMedium &&
		priority != models.PriorityHigh {
		return errors.New("invalid priority")
	}

	task.Title = title
	task.Description = description
	task.Status = status
	task.Priority = priority
	task.AssigneeID = assigneeID
	task.DueDate = dueDate
	task.UpdatedAt = time.Now().UTC()

	return s.taskRepo.UpdateTask(ctx, task)
}

func (s *TaskService) DeleteTask(
	ctx context.Context,
	taskID uuid.UUID,
	userID uuid.UUID,
) error {

	task, err := s.taskRepo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	project, err := s.projectRepo.GetProjectByID(ctx, task.ProjectID)
	if err != nil {
		return err
	}

	// 🔐 authorization
	if project.OwnerID != userID {
		return errors.New("forbidden")
	}

	return s.taskRepo.DeleteTask(ctx, taskID)
}
