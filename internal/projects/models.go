package projects

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

type userProjects struct {
	id           int64
	projectsList []project
}
