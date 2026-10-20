package auth

import "errors"

var ErrInvalidToken = errors.New("invalid token")
var ErrTokenNotFound = errors.New("token not found")
