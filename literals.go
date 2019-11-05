package pars

import "github.com/ktnyt/ascii"

// IntFunc is a parser function for matching an integer literal.
func Int(state *State, result *Result) error {
	n := 0
	s := 1

	state.Push()

	if err := state.Want(1); err != nil {
		return NewTraceError("Int", err)
	}

	if state.Head() == '-' {
		s = -1
		state.Advance()
	}

	if err := state.Want(1); err != nil {
		state.Pop()
		return NewTraceError("Int", err)
	}

	if !ascii.IsDigit(state.Head()) {
		state.Pop()
		return NewMismatchError("Int", "digit", state.Position())
	}

	// If the first byte is '0' yield 0.
	if state.Is('0') {
		state.Drop()
		state.Advance()
		result.SetValue(0)
		return nil
	}

	// The first byte is a non-zero digit
	n = n*10 + int(state.Head()-'0')
	state.Advance()
	state.Drop()

	for {
		state.Push()
		if err := state.Want(1); err != nil {
			state.Pop()
			result.SetValue(n * s)
			return nil
		}
		if !ascii.IsDigit(state.Head()) {
			state.Pop()
			result.SetValue(n * s)
			return nil
		}
		n = n*10 + int(state.Head()-'0')
		state.Drop()
		state.Advance()
	}
}
