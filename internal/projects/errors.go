package projects

import "errors"

var ErrAlreadyExists = errors.New("project with this name already exists")
var ErrInvalidProjectName = errors.New("invalid project name")
var ErrInvalidBranchName = errors.New("invalid branch name")
var ErrSizeName = errors.New("project name or branch name cannot be more than 100 characters")
var ErrProjectNotFound = errors.New("project not found")
