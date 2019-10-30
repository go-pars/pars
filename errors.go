package pars

import "fmt"

// ParserError represents a generic parser error.
type ParserError struct {
	err      error
	position Position
}

// NewParserError creates a new ParserError.
func NewParserError(v interface{}, position Position) error {
	return &ParserError{asError(v), position}
}

// Error satisfies the error interface.
func (e ParserError) Error() string {
	return fmt.Sprintf("%v at %s", e.err, e.position)
}

// MismatchError indicates that a parser has failed to match a specific
// sequence of bytes at the given position.
type MismatchError struct {
	expected string
	position Position
}

// NewMismatchError creates a new MismatchError.
func NewMismatchError(name string, expected interface{}, position Position) error {
	return NewTraceError(name, &MismatchError{asString(expected), position})
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
	return fmt.Sprintf("in `%s`: %v", e.name, e.err)
}

// Unwrap returns the underlying error value.
func (e TraceError) Unwrap() error {
	return e.err
}
