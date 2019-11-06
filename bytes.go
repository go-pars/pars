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
			if err := state.Want(1); err != nil {
				return NewTraceError("Byte()", err)
			}
			result.SetToken([]byte{state.Head()})
			state.Advance()
			return nil
		}
	case 1:
		c := p[0]
		rep := ascii.Rep(c)
		name := fmt.Sprintf("Byte(%s)", rep)

		return func(state *State, result *Result) error {
			if err := state.Want(1); err != nil {
				return NewTraceError(name, err)
			}
			if !state.Is(c) {
				return NewMismatchError(name, rep, state.Position())
			}
			result.SetToken([]byte{c})
			state.Advance()
			return nil
		}
	default:
		reps := strings.Join(ascii.Reps(p), ", ")
		name := fmt.Sprintf("Byte(%s)", reps)

		s := string(p)
		mismatch := func(c byte) bool { return strings.IndexByte(s, c) < 0 }

		return func(state *State, result *Result) error {
			if err := state.Want(1); err != nil {
				return NewTraceError(name, err)
			}
			c := state.Head()
			if mismatch(c) {
				return NewMismatchError(name, reps, state.Position())
			}
			result.SetToken([]byte{c})
			state.Advance()
			return nil
		}
	}
}

func sign(i int) int {
	if i > 0 {
		return 1
	}
	if i < 0 {
		return -1
	}
	return 0
}

// ByteRange will match any byte within the given range.
func ByteRange(begin, end byte) Parser {
	switch sign(int(end - begin)) {
	case -1:
		panic(fmt.Errorf("byte `%s` is greater than `%s`", ascii.Rep(begin), ascii.Rep(end)))
	case 0:
		return Byte(begin)
	default:
		name := fmt.Sprintf("ByteRange(%s, %s)", ascii.Rep(begin), ascii.Rep(end))
		rep := fmt.Sprintf("in range %s-%s", ascii.Rep(begin), ascii.Rep(end))

		return func(state *State, result *Result) error {
			if err := state.Want(1); err != nil {
				return NewTraceError(name, err)
			}
			c := state.Head()
			if c < begin || end < c {
				return NewMismatchError(name, rep, state.Position())
			}
			result.SetToken([]byte{c})
			state.Advance()
			return nil
		}
	}
}

// Bytes will match the given sequence of bytes.
func Bytes(p []byte) Parser {
	n := len(p)
	switch n {
	case 0:
		return Epsilon
	case 1:
		return Byte(p[0])
	default:
		reps := fmt.Sprintf("[%s]", strings.Join(ascii.Reps(p), ", "))
		name := fmt.Sprintf("Bytes([%s])", reps)

		return func(state *State, result *Result) error {
			if err := state.Want(n); err != nil {
				return NewTraceError(name, err)
			}
			if !bytes.Equal(state.Buffer(), p) {
				return NewMismatchError(name, reps, state.Position())
			}
			result.SetToken(p)
			state.Advance()
			return nil
		}
	}
}
