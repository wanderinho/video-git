package projects

type CreateProjectRequest struct {
	Name       string  `json:"name"`
	Desc       *string `json:"description"`
	BranchName string  `json:"branch_name"`
}

type CreateProjectResponse struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Desc          *string `json:"description"`
	DefaultBranch string  `json:"default_branch"`
}
