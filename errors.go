package pars

import (
	"fmt"
)

// Error is a generic parser error.
type Error struct {
	what string
	pos  Position
}

// NewError creates a new Error.
func NewError(what string, pos Position) error {
	return Error{what, pos}
}

// Error satisfies the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%s at %s", e.what, e.pos)
}
