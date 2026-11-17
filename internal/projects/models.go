package projects

import "time"

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
