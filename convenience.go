package pars

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/ktnyt/ascii"
)

func untilByte(e byte) Parser {
	name := fmt.Sprintf("Until(%s)", ascii.Rep(e))

	return func(state *State, result *Result) error {
		state.Push()

		c, err := Next(state)
		if err != nil {
			state.Pop()
			return NewNestedError(name, err)
		}

		for c != e {
			state.Advance()
			c, err = Next(state)
			if err != nil {
				state.Pop()
				return NewNestedError(name, err)
			}
		}

		p, err := Trail(state)
		if err != nil {
			return NewNestedError(name, err)
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
		name := fmt.Sprintf("Until(%s)", strings.Join(ascii.Reps(e), ", "))

		return func(state *State, result *Result) error {
			state.Push()

			for {
				if err := state.Request(len(e)); err != nil {
					state.Pop()
					return NewNestedError(name, err)
				}

				if bytes.Equal(state.Buffer(), e) {
					p, err := Trail(state)
					if err != nil {
						return NewNestedError(name, err)
					}
					result.SetToken(p)
					return nil
				}

				if err := Skip(state, 1); err != nil {
					state.Pop()
					return NewNestedError(name, err)
				}
			}
		}
	}
}

func untilFilter(filter ascii.Filter) Parser {
	v := reflect.ValueOf(filter)
	f := runtime.FuncForPC(v.Pointer())
	name := fmt.Sprintf("Until(%s)", f.Name())

	return func(state *State, result *Result) error {
		state.Push()

		c, err := Next(state)
		if err != nil {
			state.Pop()
			return NewNestedError(name, err)
		}

		for !filter(c) {
			state.Advance()
			c, err = Next(state)
			if err != nil {
				state.Pop()
				return NewNestedError(name, err)
			}
		}

		p, err := Trail(state)
		if err != nil {
			return NewNestedError(name, err)
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
					return NewNestedError("Until", err)
				}
				state.Push()
			}
			state.Pop()

			p, err := Trail(state)
			if err != nil {
				return NewNestedError("Until", err)
			}
			result.SetToken(p)
			return nil
		}
	}
}

// EOL matches an end of a line (a newline byte or the end of state).
func EOL(state *State, result *Result) error {
	c, err := Next(state)
	if err != nil {
		return nil
	}
	if c != '\n' {
		return NewError("expected newline or end of state", state.Position())
	}
	state.Advance()
	return nil
}

// Line matches up to a newline byte or the end of state.
func Line(state *State, result *Result) error {
	state.Push()
	c, err := Next(state)
	for err == nil && c != '\n' {
		state.Advance()
		c, err = Next(state)
	}
	p, err := Trail(state)
	if err != nil && err != io.EOF {
		panic(err)
	}
	result.SetToken(p)
	Skip(state, 1)
	return nil
}
