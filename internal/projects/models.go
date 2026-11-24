package projects

import (
	"time"
)

type createProjectParams struct {
	ownerID     int64
	name        string
	description *string
	branchName  string
}

type project struct {
	id            int64
	name          string
	desc          *string
	defaultBranch string
}

type fullProject struct {
	id                int64
	ownerID           int64
	name              string
	desc              *string
	createdAt         time.Time
	defaultBranchID   int64
	defaultBranchName string
}

type userProjects struct {
	id           int64
	projectsList []project
}

type createBranchRequest struct {
	userID    int64
	commitID  int64
	projectID int64
	name      string
}

type createBranchResponse struct {
	id           int64
	name         string
	headCommitID int64
	createdAt    time.Time
}

type branchInfo struct {
	id           int64
	name         string
	createdBy    string
	createdAt    time.Time
	headCommitID *int64
}

type projectBranches struct {
	projectID       int64
	defaultBranchID int64
	branches        []branchInfo
}
