package pars

import "strings"

// Byte will match the given byte.
func Byte(c byte) Parser {
	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return err
		}
		if state.Head() != c {
			return NewMismatchError("Byte", c, state.Position())
		}
		result.SetToken(asBytes(c))
		state.Advance()
		return nil
	}
}

// AnyByte will match any of the given bytes.
func AnyByte(p ...byte) Parser {
	s := string(p)
	mismatch := func(c byte) bool { return strings.IndexByte(s, c) < 0 }

	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return err
		}
		c := state.Head()
		if mismatch(c) {
			return NewMismatchError("Byte", s, state.Position())
		}
		result.SetToken(asBytes(c))
		state.Advance()
		return nil
	}
}
