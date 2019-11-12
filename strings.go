package pars

import (
	"bytes"
	"fmt"
)

func String(s string) Parser {
	what := fmt.Sprintf("expected \"%s\"", s)
	p := []byte(s)

	return func(state *State, result *Result) error {
		if err := state.Request(len(p)); err != nil {
			return err
		}
		if !bytes.Equal(state.Buffer(), p) {
			return NewError(what, state.Position())
		}
		result.SetValue(s)
		state.Advance()
		return nil
	}
}
