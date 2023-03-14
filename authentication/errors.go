package authentication

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUnauthorized = errors.New("unauthorized")
var ErrUserNotFound = errors.New("user not found")
