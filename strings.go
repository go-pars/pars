package pars

import (
	"bytes"
	"fmt"
)

func String(s string) Parser {
	name := fmt.Sprintf("String(\"%s\")", s)
	p := []byte(s)

	return func(state *State, result *Result) error {
		if err := state.Want(len(p)); err != nil {
			return NewTraceError(name, err)
		}
		if !bytes.Equal(state.Buffer(), p) {
			return NewMismatchError(name, s, state.Position())
		}
		result.SetValue(s)
		state.Advance()
		return nil
	}
}
