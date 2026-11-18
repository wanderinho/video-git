package projects

import (
	"context"
	"strings"
	"unicode/utf8"
)

type ProjectService struct {
	projectRepo Repository
}

type Repository interface {
	Create(ctx context.Context, createProjectForm createProjectParams) (project, error)
	ListByUser(ctx context.Context, userID int64) (userProjects, error)
	GetProjectByID(ctx context.Context, projectID int64, userID int64) (fullProject, error)
	AddMember(ctx context.Context, projectID, ownerID, memberID int64) error
}

func NewProjectService(projectRepo Repository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

const maxNameLength = 100

func (ps *ProjectService) Create(ctx context.Context, userID int64, desc *string, name, branchName string) (project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return project{}, ErrInvalidProjectName
	}

	branchName = strings.TrimSpace(branchName)
	if branchName == "" {
		return project{}, ErrInvalidBranchName
	}

	if utf8.RuneCountInString(name) > maxNameLength || utf8.RuneCountInString(branchName) > maxNameLength {
		return project{}, ErrSizeName
	}

	params := createProjectParams{
		ownerID:     userID,
		name:        name,
		description: desc,
		branchName:  branchName,
	}

	newProject, err := ps.projectRepo.Create(ctx, params)
	if err != nil {
		return project{}, err
	}

	return newProject, nil
}

func (ps *ProjectService) ListByUser(ctx context.Context, userID int64) (userProjects, error) {
	return ps.projectRepo.ListByUser(ctx, userID)
}

func (ps *ProjectService) GetProjectByID(ctx context.Context, projectID, userID int64) (fullProject, error) {
	return ps.projectRepo.GetProjectByID(ctx, projectID, userID)
}

func (ps *ProjectService) AddMember(ctx context.Context, projectID, ownerID, memberID int64) error {
	return ps.projectRepo.AddMember(ctx, projectID, ownerID, memberID)
}