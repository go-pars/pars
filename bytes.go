package pars

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ktnyt/ascii"
)

// Byte will attempt to match the next single byte.
// If no bytes are given, it will match any byte.
// Otherwise, the given bytes will be tested for a match.
func Byte(p ...byte) Parser {
	switch len(p) {
	case 0:
		return func(state *State, result *Result) error {
			if err := state.Request(1); err != nil {
				return err
			}
			result.SetToken([]byte{state.Buffer()[0]})
			state.Advance()
			return nil
		}
	case 1:
		e := p[0]
		what := fmt.Sprintf("expected `%s`", ascii.Rep(e))

		return func(state *State, result *Result) error {
			c, err := Next(state)
			if err != nil {
				return err
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
		what := fmt.Sprintf("expected one of [%s]", reps)

		s := string(p)
		mismatch := func(c byte) bool { return strings.IndexByte(s, c) < 0 }

		return func(state *State, result *Result) error {
			c, err := Next(state)
			if err != nil {
				return err
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

// ByteRange will match any byte within the given range.
func ByteRange(begin, end byte) Parser {
	if begin < end {
		what := fmt.Sprintf("expected in range %s-%s", ascii.Rep(begin), ascii.Rep(end))

		return func(state *State, result *Result) error {
			c, err := Next(state)
			if err != nil {
				return err
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

// Bytes will match the given sequence of bytes.
func Bytes(p []byte) Parser {
	if n := len(p); n > 0 {
		reps := fmt.Sprintf("[%s]", strings.Join(ascii.Reps(p), ", "))
		what := fmt.Sprintf("expected [%s]", reps)

		return func(state *State, result *Result) error {
			if err := state.Request(n); err != nil {
				return err
			}
			if !bytes.Equal(state.Buffer(), p) {
				return NewError(what, state.Position())
			}
			result.SetToken(p)
			state.Advance()
			return nil
		}
	}
	panic("no bytes given")
}
