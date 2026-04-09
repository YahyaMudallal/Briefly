package apperrors

import "errors"

// errors definitions
var (
	ErrValidation = errors.New("validation failed")
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
	ErrInternal = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
)