package pars

import "fmt"

// ParserError represents a generic parser error.
type ParserError struct {
	message  string
	position int
}

// NewParserError creates a new ParserError.
func NewParserError(message string, position int) error {
	return &ParserError{message: message, position: position}
}

// Error satisfies the error interface
func (e *ParserError) Error() string {
	return fmt.Sprintf("%s at position %d", e.message, e.position)
}

// MismatchError represents a parser mismatch.
type MismatchError struct {
	expected []byte
	position int
	parser   string
	not      bool
}

// NewMismatchError creates a new MismatchError.
func NewMismatchError(parser string, expected []byte, position int) error {
	return &MismatchError{
		expected: expected,
		position: position,
		parser:   parser,
		not:      false,
	}
}

// NewNotMismatchError creates a new MismatchError.
func NewNotMismatchError(parser string, expected []byte, position int) error {
	return &MismatchError{
		parser:   parser,
		expected: expected,
		position: position,
		not:      true,
	}
}

// Error satisfies the error interface.
func (e *MismatchError) Error() string {
	if e.not {
		return fmt.Sprintf("`%s` expected not `%s` at position %d", e.parser, e.expected, e.position)
	}
	return fmt.Sprintf("`%s` expected `%s` at position %d", e.parser, e.expected, e.position)
}

// TraceError traces nested errors.
type TraceError struct {
	parser string
	err    error
}

// NewTraceError creates a new TraceError
func NewTraceError(parser string, err error) error {
	return &TraceError{parser: parser, err: err}
}

// Error satisfies the error interface.
func (e *TraceError) Error() string {
	return fmt.Sprintf("in parser `%s`:\n%s", e.parser, e.err.Error())
}
