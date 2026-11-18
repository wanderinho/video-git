package projects

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: pool}
}

func (pr *ProjectRepository) Create(ctx context.Context, createProjectForm createProjectParams) (project, error) {
	tx, err := pr.db.Begin(ctx)
	if err != nil {
		return project{}, fmt.Errorf("begin create project transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var projectID int64
	err = tx.QueryRow(ctx, "INSERT INTO projects (owner_id, name, description) VALUES ($1, $2, $3) RETURNING id", createProjectForm.ownerID, createProjectForm.name, createProjectForm.description).Scan(&projectID)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) &&
			pgError.Code == "23505" &&
			pgError.ConstraintName == "projects_owner_name_unique" {
			return project{}, ErrAlreadyExists
		}

		return project{}, fmt.Errorf("insert project: %w", err)
	}

	var branchID int64
	err = tx.QueryRow(ctx, "INSERT INTO branches (project_id, name, created_by) VALUES ($1, $2, $3) RETURNING id", projectID, createProjectForm.branchName, createProjectForm.ownerID).Scan(&branchID)
	if err != nil {
		return project{}, fmt.Errorf("insert default branch: %w", err)
	}

	tag, err := tx.Exec(ctx, "UPDATE projects SET default_branch_id = $1 WHERE id = $2", branchID, projectID)
	if err != nil {
		return project{}, fmt.Errorf("update default branch id: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return project{}, errors.New("project update affected unexpected number of rows")
	}

	_, err = tx.Exec(ctx, `INSERT INTO project_members (project_id, user_id) VALUES ($1, $2)`, projectID, createProjectForm.ownerID)
	if err != nil {
		return project{}, fmt.Errorf("add project owner as member: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return project{}, fmt.Errorf("commit create project transaction: %w", err)
	}

	return project{
		id:            projectID,
		name:          createProjectForm.name,
		desc:          createProjectForm.description,
		defaultBranch: createProjectForm.branchName,
	}, nil

}

func (pr *ProjectRepository) ListByUser(ctx context.Context, userID int64) (userProjects, error) {
	const query = `
        SELECT p.id, p.name, p.description, b.name
        FROM project_members pm
        JOIN projects p ON p.id = pm.project_id
        JOIN branches b ON b.id = p.default_branch_id
        WHERE pm.user_id = $1
        ORDER BY p.created_at DESC, p.id DESC
    `

	rows, err := pr.db.Query(ctx, query, userID)
	if err != nil {
		return userProjects{}, fmt.Errorf("get user projects: %w", err)
	}
	defer rows.Close()

	projects := userProjects{
		id:           userID,
		projectsList: make([]project, 0),
	}

	for rows.Next() {
		var project project

		err := rows.Scan(
			&project.id,
			&project.name,
			&project.desc,
			&project.defaultBranch,
		)
		if err != nil {
			return userProjects{}, fmt.Errorf("scan project: %w", err)
		}

		projects.projectsList = append(projects.projectsList, project)
	}

	if err := rows.Err(); err != nil {
		return userProjects{}, fmt.Errorf("iterate projects: %w", err)
	}

	return projects, nil
}

func (pr *ProjectRepository) GetProjectByID(ctx context.Context, projectID, userID int64) (fullProject, error) {
	const query = `
        SELECT p.id, p.owner_id, p.name, p.description, p.created_at, b.id, b.name
        FROM projects p
        JOIN project_members pm ON pm.project_id = p.id
        JOIN branches b ON b.id = p.default_branch_id
        WHERE p.id = $1 AND pm.user_id = $2
    `

	var project fullProject
	err := pr.db.QueryRow(ctx, query, projectID, userID).Scan(&project.id, &project.ownerID, &project.name, &project.desc, &project.createdAt, &project.defaultBranchID, &project.defaultBranchName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fullProject{}, ErrProjectNotFound
		}
		return fullProject{}, fmt.Errorf("get project by id: %w", err)
	}

	return project, nil
}

func (pr *ProjectRepository) AddMember(ctx context.Context, projectID, ownerID, memberID int64) error {
	const query = `
        INSERT INTO project_members (project_id, user_id)
        SELECT id, $3
        FROM projects
        WHERE id = $1 AND owner_id = $2
    `

	tag, err := pr.db.Exec(ctx, query, projectID, ownerID, memberID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505" &&
				pgErr.ConstraintName == "project_members_pkey":
				return ErrMemberAlreadyExists

			case pgErr.Code == "23503" &&
				pgErr.ConstraintName == "project_members_user_id_fkey":
				return ErrUserNotFound
			}
		}

		return fmt.Errorf("add project member: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrProjectNotFound
	}

	return nil
}

func (pr *ProjectRepository) ChangeDefaultBranch(ctx context.Context, userID, projectID, branchID int64) error {
	tx, err := pr.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin change default branch transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var id int64
	err = tx.QueryRow(ctx, `SELECT id FROM projects WHERE id = $1 AND owner_id = $2 FOR UPDATE`, projectID, userID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectNotFound
		}
		return fmt.Errorf("check project ownership: %w", err)
	}

	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM branches WHERE id = $1 AND project_id = $2)`, branchID, projectID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check branch: %w", err)
	}

	if !exists {
		return ErrBranchNotFound
	}

	_, err = tx.Exec(ctx, `UPDATE projects SET default_branch_id = $1 WHERE id = $2`, branchID, projectID)
	if err != nil {
		return fmt.Errorf("update default branch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit change default branch: %w", err)
	}

	return nil
}
