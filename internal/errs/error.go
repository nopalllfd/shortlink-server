package errs

import "errors"

var (
	ErrInternalServer  = errors.New("internal server error")
	ErrExistingEmail   = errors.New("email already registered")
	ErrUserNotFound    = errors.New("user not found")
	InvalidCredentials = errors.New("invalid credentials")

	ErrSlugAlreadyExists = errors.New("slug already exists")

	ErrLinkNotFound = errors.New("link not found")

	ErrCannotUserReserveWord = errors.New("cannot use reserved words")
	ErrSlugNotFound          = errors.New("slug not found")

	ErrMinimumSlug = errors.New("minimum 6 chars for custom slug")
)
