package pars

import (
	"bytes"
	"fmt"
	"strings"

	ascii "gopkg.in/ascii.v1"
)

// Byte creates a Parser which will attempt to match the next single byte.
// If no bytes are given, it will match any byte.
// Otherwise, the given bytes will be tested for a match.
func Byte(p ...byte) Parser {
	switch len(p) {
	case 0:
		return func(state *State, result *Result) error {
			if err := state.Request(1); err != nil {
				return NewNestedError("Byte", err)
			}
			result.SetToken([]byte{state.Buffer()[0]})
			state.Advance()
			return nil
		}
	case 1:
		e := p[0]
		rep := ascii.Rep(e)
		name := fmt.Sprintf("Byte(%s)", rep)
		what := fmt.Sprintf("expected `%s`", rep)

		return func(state *State, result *Result) error {
			c, err := Next(state)
			if err != nil {
				return NewNestedError(name, err)
			}
			if c != e {
				return NewError(what, state.Position())
			}
			result.SetToken([]byte{c})
			state.Advance()
			return nil
		}
	default:
		reps := strings.Join(ascii.Reps(p), ", ")
		name := fmt.Sprintf("Byte(%s)", reps)
		what := fmt.Sprintf("expected one of [%s]", reps)

		s := string(p)
		mismatch := func(c byte) bool { return strings.IndexByte(s, c) < 0 }

		return func(state *State, result *Result) error {
			c, err := Next(state)
			if err != nil {
				return NewNestedError(name, err)
			}
			if mismatch(c) {
				return NewError(what, state.Position())
			}
			result.SetToken([]byte{c})
			state.Advance()
			return nil
		}
	}
}

// ByteRange creates a Parser which will attempt to match a byte between the
// given range inclusively.
func ByteRange(begin, end byte) Parser {
	if begin < end {
		rbegin, rend := ascii.Rep(begin), ascii.Rep(end)
		name := fmt.Sprintf("ByteRange(%s, %s)", rbegin, rend)
		what := fmt.Sprintf("expected in range %s-%s", rbegin, rend)

		return func(state *State, result *Result) error {
			c, err := Next(state)
			if err != nil {
				return NewNestedError(name, err)
			}
			if c < begin || end < c {
				return NewError(what, state.Position())
			}
			result.SetToken([]byte{c})
			state.Advance()
			return nil
		}
	}
	panic("invalid byte range")
}

// Bytes creates a Parser which will attempt to match the given sequence of bytes.
func Bytes(p []byte) Parser {
	reps := fmt.Sprintf("[%s]", strings.Join(ascii.Reps(p), ", "))
	name := fmt.Sprintf("Bytes([%s])", reps)
	what := fmt.Sprintf("expected [%s]", reps)

	return func(state *State, result *Result) error {
		if err := state.Request(len(p)); err != nil {
			return NewNestedError(name, err)
		}
		if !bytes.Equal(state.Buffer(), p) {
			return NewError(what, state.Position())
		}
		result.SetToken(p)
		state.Advance()
		return nil
	}
}
