package repository

import (
	"context"
	"errors"
	"taskflow/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}
func (r *ProjectRepository) CreateProject(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO projects (id, name, description, owner_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		project.ID,
		project.Name,
		project.Description,
		project.OwnerID,
		project.CreatedAt,
	)

	return err
}

func (r *ProjectRepository) GetProjectsByUser(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
	offset int,
) ([]models.Project, error) {

	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE p.owner_id = $1 OR t.assignee_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.OwnerID,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	query := `
		SELECT id, name, description, owner_id, created_at
		FROM projects
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var p models.Project

	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.OwnerID,
		&p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, project *models.Project) error {
	query := `
		UPDATE projects
		SET name = $1, description = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query,
		project.Name,
		project.Description,
		project.ID,
	)

	return err
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
