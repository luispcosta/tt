package utils

import "fmt"

// NotFoundError represents the specific error NotFound (useful when something cannot be found, like an activity)
type NotFoundError struct {
	Err string
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(err string) *NotFoundError {
	return &NotFoundError{Err: err}
}

// Error implements the Error method from the error interface
func (err *NotFoundError) Error() string {
	return fmt.Sprintf("NotFound Error: %s", err.Err)
}
