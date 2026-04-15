package seed

import (
	"context"
	"log"
	"time"

	"taskflow/internal/db"
	"taskflow/internal/models"
	"taskflow/internal/repository"
	"taskflow/internal/utils"

	"github.com/google/uuid"
)

func Run(database *db.Database) {
	ctx := context.Background()

	userRepo := repository.NewUserRepository(database.Pool)
	projectRepo := repository.NewProjectRepository(database.Pool)
	taskRepo := repository.NewTaskRepository(database.Pool)

	password, _ := utils.HashPassword("password123")

	user := &models.User{
		ID:        uuid.New(),
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  password,
		CreatedAt: time.Now(),
	}

	err := userRepo.CreateUser(ctx, user)
	if err != nil {
		log.Println("User already exists, skipping seed")
	}

	project := &models.Project{
		ID:          uuid.New(),
		Name:        "Demo Project",
		Description: "Seeded project",
		OwnerID:     user.ID,
		CreatedAt:   time.Now(),
	}

	_ = projectRepo.CreateProject(ctx, project)

	task1 := &models.Task{
		ID:        uuid.New(),
		Title:     "Task 1",
		Status:    models.StatusTodo,
		Priority:  models.PriorityHigh,
		ProjectID: project.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	task2 := &models.Task{
		ID:        uuid.New(),
		Title:     "Task 2",
		Status:    models.StatusInProgress,
		Priority:  models.PriorityMedium,
		ProjectID: project.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	task3 := &models.Task{
		ID:        uuid.New(),
		Title:     "Task 3",
		Status:    models.StatusDone,
		Priority:  models.PriorityLow,
		ProjectID: project.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_ = taskRepo.CreateTask(ctx, task1)
	_ = taskRepo.CreateTask(ctx, task2)
	_ = taskRepo.CreateTask(ctx, task3)

	log.Println("✅ Seed data inserted")
}
