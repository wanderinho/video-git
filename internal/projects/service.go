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
	Create(ctx context.Context, createProjectForm CreateProjectParams) (CreatedProject, error)
}

func NewProjectService(projectRepo Repository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

const maxNameLength = 100

func (ps *ProjectService) Create(ctx context.Context, userID int64, desc *string, name, branchName string) (CreatedProject, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return CreatedProject{}, ErrInvalidProjectName
	}

	branchName = strings.TrimSpace(branchName)
	if branchName == "" {
		return CreatedProject{}, ErrInvalidBranchName
	}

	if utf8.RuneCountInString(name) > maxNameLength || utf8.RuneCountInString(branchName) > maxNameLength {
		return CreatedProject{}, ErrSizeName
	}

	params := CreateProjectParams{
		OwnerID:     userID,
		Name:        name,
		Description: desc,
		BranchName:  branchName,
	}

	newProject, err := ps.projectRepo.Create(ctx, params)
	if err != nil {
		return CreatedProject{}, err
	}

	return newProject, nil
}
