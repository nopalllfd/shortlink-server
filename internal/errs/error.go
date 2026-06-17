package errs

import "errors"

var (
	ErrInternalServer  = errors.New("internal server error")
	ErrExistingEmail   = errors.New("email already registered")
	ErrUserNotFound    = errors.New("user not found")
	InvalidCredentials = errors.New("invalid credentials")
)
