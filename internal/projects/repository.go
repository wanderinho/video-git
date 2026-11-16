package projects

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: pool}
}

func (pr *ProjectRepository) Create(ctx context.Context, createProjectForm CreateProjectParams) (CreatedProject, error) {
	tx, err := pr.db.Begin(ctx)
	if err != nil {
		return CreatedProject{}, fmt.Errorf("begin create project transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var projectID int64
	err = tx.QueryRow(ctx, "INSERT INTO projects (owner_id, name, description) VALUES ($1, $2, $3) RETURNING id", createProjectForm.OwnerID, createProjectForm.Name, createProjectForm.Description).Scan(&projectID)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) &&
			pgError.Code == "23505" &&
			pgError.ConstraintName == "projects_owner_name_unique" {
			return CreatedProject{}, ErrAlreadyExists
		}

		return CreatedProject{}, fmt.Errorf("insert project: %w", err)
	}

	var branchID int64
	err = tx.QueryRow(ctx, "INSERT INTO branches (project_id, name, created_by) VALUES ($1, $2, $3) RETURNING id", projectID, createProjectForm.BranchName, createProjectForm.OwnerID).Scan(&branchID)
	if err != nil {
		return CreatedProject{}, fmt.Errorf("insert default branch: %w", err)
	}

	tag, err := tx.Exec(ctx, "UPDATE projects SET default_branch_id = $1 WHERE id = $2", branchID, projectID)
	if err != nil {
		return CreatedProject{}, fmt.Errorf("update default branch id: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return CreatedProject{}, errors.New("project update affected unexpected number of rows")
	}

	if err := tx.Commit(ctx); err != nil {
		return CreatedProject{}, fmt.Errorf("commit create project transaction: %w", err)
	}

	return CreatedProject{
		ID:            projectID,
		Name:          createProjectForm.Name,
		Desc:          createProjectForm.Description,
		DefaultBranch: createProjectForm.BranchName,
	}, nil

}
