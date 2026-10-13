package users

import "errors"

var ErrAlreadyExists = errors.New("username or email already exists")
var ErrLongPassword = errors.New("password too long")
var ErrInvalidCredentials = errors.New("incorrect username or password")
