package pars

import (
	"errors"
	"fmt"
)

var errNoChildren = errors.New("result does not have children")

// Error is a generic parser error.
type Error struct {
	what string
	pos  Position
}

// NewError creates a new Error.
func NewError(what string, pos Position) error { return Error{what, pos} }

// Error satisfies the error interface.
func (e Error) Error() string { return fmt.Sprintf("%s at %s", e.what, e.pos) }

// NestedError is a nested error type.
type NestedError struct {
	name string
	err  error
}

// NewNestedError creates a new NestedError.
func NewNestedError(name string, err error) error {
	return NestedError{name, err}
}

// Error satisfies the error interface.
func (e NestedError) Error() string {
	return fmt.Sprintf("in %s:\n%s", e.name, e.err)
}

// Unwrap returns the internal error value.
func (e NestedError) Unwrap() error { return e.err }
