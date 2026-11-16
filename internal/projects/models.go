package projects

type CreateProjectParams struct {
	OwnerID     int64
	Name        string
	Description *string
	BranchName  string
}

type CreatedProject struct {
	ID int64
	Name string
	Desc *string
	DefaultBranch string
}