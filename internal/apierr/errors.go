package apierr

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrBadRequest  = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden   = errors.New("forbidden")
	ErrNotFound    = errors.New("not found")
	ErrConflict    = errors.New("conflict")
)

func BadRequest(message string) error {
	return fmt.Errorf("%w: %s", ErrBadRequest, message)
}

func Unauthorized(message string) error {
	return fmt.Errorf("%w: %s", ErrUnauthorized, message)
}

func Forbidden(message string) error {
	return fmt.Errorf("%w: %s", ErrForbidden, message)
}

func NotFound(message string) error {
	return fmt.Errorf("%w: %s", ErrNotFound, message)
}

func Conflict(message string) error {
	return fmt.Errorf("%w: %s", ErrConflict, message)
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	prefixes := []string{
		"bad request: ",
		"unauthorized: ",
		"forbidden: ",
		"not found: ",
		"conflict: ",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(message, prefix) {
			return strings.TrimPrefix(message, prefix)
		}
	}
	return message
}
