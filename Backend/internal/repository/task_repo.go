package repository

import (
	"context"
	"errors"
	"fmt"
	"taskflow/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO tasks 
		(id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`

	_, err := r.db.Exec(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.ProjectID,
		task.AssigneeID,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return err
}

func (r *TaskRepository) GetTasksByProject(
	ctx context.Context,
	projectID uuid.UUID,
	status *models.TaskStatus,
	assigneeID *uuid.UUID,
) ([]models.Task, error) {

	baseQuery := `
		SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
		FROM tasks
		WHERE project_id = $1
	`

	args := []interface{}{projectID}
	argIndex := 2

	// dynamic filters
	if status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *status)
		argIndex++
	}

	if assigneeID != nil {
		baseQuery += fmt.Sprintf(" AND assignee_id = $%d", argIndex)
		args = append(args, *assigneeID)
		argIndex++
	}

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Priority,
			&t.ProjectID,
			&t.AssigneeID,
			&t.DueDate,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}
func (r *TaskRepository) GetTaskByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var task models.Task

	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.ProjectID,
		&task.AssigneeID,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks
		SET title = $1,
		    description = $2,
		    status = $3,
		    priority = $4,
		    assignee_id = $5,
		    due_date = $6,
		    updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.Exec(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.AssigneeID,
		task.DueDate,
		task.UpdatedAt,
		task.ID,
	)

	return err
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
