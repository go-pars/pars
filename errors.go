package pars

import (
	"fmt"
)

// ParserError represents a generic parser error.
type ParserError struct {
	what     string
	position Position
}

// NewParserError creates a new ParserError.
func NewParserError(what string, position Position) error {
	return &ParserError{what, position}
}

// Error satisfies the error interface.
func (e ParserError) Error() string {
	return fmt.Sprintf("%s at %s", e.what, e.position)
}

// MismatchError indicates that a parser has failed to match a specific
// sequence of bytes at the given position.
type MismatchError struct {
	expected string
	position Position
}

// NewMismatchError creates a new MismatchError.
func NewMismatchError(name string, expected string, position Position) error {
	return NewTraceError(name, &MismatchError{expected, position})
}

// Error satisfies the error interface.
func (e MismatchError) Error() string {
	return fmt.Sprintf("expected `%s` at %s", e.expected, e.position)
}

// TraceError holds nested error values.
type TraceError struct {
	name string
	err  error
}

// NewTraceError creates a new trace error.
func NewTraceError(name string, err error) error {
	return &TraceError{name, err}
}

// Error satisfies the error interface.
func (e TraceError) Error() string {
	return fmt.Sprintf("in `%s`:\n%v", e.name, e.err)
}

// Unwrap returns the underlying error value.
func (e TraceError) Unwrap() error {
	return e.err
}
