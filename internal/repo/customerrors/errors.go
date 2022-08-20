package customerrors

import "github.com/pkg/errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrTimeout           = errors.New("deadline exceeded")
	ErrUnexpected        = errors.New("unexpected error")
)
