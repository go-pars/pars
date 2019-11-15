package pars

import (
	"bytes"
	"unicode/utf8"

	"github.com/ktnyt/ascii"
)

func untilByte(e byte) Parser {
	return func(state *State, result *Result) error {
		state.Push()

		c, err := Next(state)
		if err != nil {
			state.Pop()
			return err
		}

		for c != e {
			state.Advance()
			c, err = Next(state)
			if err != nil {
				state.Pop()
				return err
			}
		}

		p, err := Trail(state)
		if err != nil {
			return err
		}
		result.SetToken(p)
		return nil
	}
}

func untilBytes(e []byte) Parser {
	switch len(e) {
	case 0:
		panic("no bytes given to Until")
	case 1:
		return untilByte(e[0])
	default:
		return func(state *State, result *Result) error {
			state.Push()

			for {
				if err := state.Request(len(e)); err != nil {
					state.Pop()
					return err
				}

				if bytes.Equal(state.Buffer(), e) {
					p, err := Trail(state)
					if err != nil {
						return err
					}
					result.SetToken(p)
					return nil
				}

				if err := Skip(state, 1); err != nil {
					state.Pop()
					return err
				}
			}
		}
	}
}

func untilFilter(filter ascii.Filter) Parser {
	return func(state *State, result *Result) error {
		state.Push()

		c, err := Next(state)
		if err != nil {
			state.Pop()
			return err
		}

		for !filter(c) {
			state.Advance()
			c, err = Next(state)
			if err != nil {
				state.Pop()
				return err
			}
		}

		p, err := Trail(state)
		if err != nil {
			return err
		}
		result.SetToken(p)
		return nil
	}
}

// Until creates a Parser which will advance the state until the given Parser
// matches, and return extent of the buffer up to the matching state position.
func Until(q interface{}) Parser {
	switch v := q.(type) {
	case byte:
		return untilByte(v)
	case []byte:
		return untilBytes(v)
	case rune:
		p := []byte{0, 0, 0, 0}
		n := utf8.EncodeRune(p, v)
		return untilBytes(p[:n])
	case []rune:
		p := []byte(string(v))
		return untilBytes(p)
	case ascii.Filter:
		return untilFilter(v)
	default:
		p := AsParser(q)

		return func(state *State, result *Result) error {
			// Backtrack point for later.
			state.Push()

			// Start scanning.
			state.Push()
			for p(state, result) != nil {
				state.Drop()
				if err := Skip(state, 1); err != nil {
					state.Pop()
					return err
				}
				state.Push()
			}
			state.Pop()

			p, err := Trail(state)
			if err != nil {
				return err
			}
			result.SetToken(p)
			return nil
		}
	}
}

// Line matches up to a newline byte or the end of state.
func Line(state *State, result *Result) error {
	state.Push()
	c, err := Next(state)
	for err != nil && c != '\n' {
		state.Advance()
		c, err = Next(state)
	}
	p, err := Trail(state)
	if err != nil {
		panic(err)
	}
	result.SetToken(p)
	if err := Skip(state, 1); err != nil {
		return err
	}
	return nil
}
