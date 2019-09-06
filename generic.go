package pars

import (
	"io"
)

// Line matches a single line of text.
func Line(state *State, result *Result) error {
	start := state.Index

	state.Mark()

	for {
		if err := state.Want(1); err != nil {
			if err == io.EOF {
				result.Value = string(state.Buffer[start:state.Index])
				return nil
			}

			state.Jump()
			return err
		}

		if state.Buffer[state.Index] == '\n' {
			result.Value = string(state.Buffer[start:state.Index])
			state.Advance(1)
			return nil
		}

		state.Advance(1)
	}
}

// Word matches a group of non-whitespace characters.
func Word(state *State, result *Result) error {
	start := state.Index

	state.Mark()

	for {
		if err := state.Want(1); err != nil {
			if err == io.EOF {
				result.Value = string(state.Buffer[start:state.Index])
				return nil
			}

			state.Jump()
			return err
		}

		if isWhitespace(state.Buffer[state.Index]) {
			result.Value = string(state.Buffer[start:state.Index])
			return nil
		}

		state.Advance(1)
	}
}
