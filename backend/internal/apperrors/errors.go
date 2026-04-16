package apperrors

import (
	"errors"
	"net/http"
)

// errors definitions
var (
	ErrValidation = errors.New("validation failed")
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
	ErrInternal = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized")
)

// FilterError take an apperror and return the corresponding HTTP status code
func FilterError(err error) int {
	switch {
	case errors.Is(err, ErrValidation):
		return http.StatusBadRequest
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}