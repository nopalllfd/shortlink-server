package errs

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrExistingEmail  = errors.New("email already registered")
)
