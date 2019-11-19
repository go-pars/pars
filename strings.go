package pars

import (
	"bytes"
	"fmt"
)

// String creates a Parser which will attempt to match the given string.
func String(s string) Parser {
	name := fmt.Sprintf(`String(%s)`, s)
	what := fmt.Sprintf(`expected "%s"`, s)
	p := []byte(s)

	return func(state *State, result *Result) error {
		if err := state.Request(len(p)); err != nil {
			return NewNestedError(name, err)
		}
		if !bytes.Equal(state.Buffer(), p) {
			return NewError(what, state.Position())
		}
		result.SetValue(s)
		state.Advance()
		return nil
	}
}
