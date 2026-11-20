package projects

import "time"

type CreateProjectRequest struct {
	Name       string  `json:"name"`
	Desc       *string `json:"description"`
	BranchName string  `json:"branch_name"`
}

type ProjectResponse struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Desc          *string `json:"description"`
	DefaultBranch string  `json:"default_branch"`
}

type ProjectsList struct {
	UserID   int64             `json:"user_id"`
	Projects []ProjectResponse `json:"projects"`
}

type FullProjectResponse struct {
	ID                int64     `json:"id"`
	OwnerID           int64     `json:"owner_id"`
	Name              string    `json:"name"`
	Desc              *string   `json:"description"`
	CreatedAt         time.Time `json:"created_at"`
	DefaultBranchID   int64     `json:"default_branch_id"`
	DefaultBranchName string    `json:"default_branch_name"`
}

type AddMemberRequest struct {
	UserID int64 `json:"user_id"`
}

type ChangeBranchRequest struct {
	BranchID int64 `json:"branch_id"`
}

type CreateBranchRequest struct {
	CommitID int64  `json:"commit_id"`
	Name     string `json:"name"`
}

type CreateBranchResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	HeadCommitID int64     `json:"head_commit_id"`
	CreatedAt    time.Time `json:"created_at"`
}
