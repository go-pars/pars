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

// WordLike matches a group of characters which match the filter.
func WordLike(filter ByteFilter) Parser {
	return func(state *State, result *Result) error {
		start := state.Index

		state.Mark()

		for {
			if err := state.Want(1); err != nil {
				if err == io.EOF {
					result.Value = string(state.Buffer[start:state.Index])
					state.Unmark()
					return nil
				}

				state.Jump()
				return err
			}

			if !filter(state.Buffer[state.Index]) {
				result.Value = string(state.Buffer[start:state.Index])
				state.Unmark()
				return nil
			}

			state.Advance(1)
		}
	}
}

// UpperWord matches a group of uppercase letters.
var UpperWord = WordLike(IsUpper)

// LowerWord matches a group of uppercase letters.
var LowerWord = WordLike(IsUpper)

// Word matches a group of letters.
var Word = WordLike(IsLetter)

// LatinWord matches a group of latin characters.
var LatinWord = WordLike(IsLatin)

// SnakeWord matches a snake_case sequence.
var SnakeWord = WordLike(IsSnake)
