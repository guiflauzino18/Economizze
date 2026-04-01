package errors

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

// Erros de Dominio
var (
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidOperation   = errors.New("invalid operation")
	ErrInsufficientFounds = errors.New("insufficient founds")
)

// Sobrescreve método Error
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %q: %s", e.Field, e.Message)
}

func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
