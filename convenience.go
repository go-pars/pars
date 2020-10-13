package pars

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/go-ascii/ascii"
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

		p, _ := Trail(state)
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
					p, _ := Trail(state)
					result.SetToken(p)
					return nil
				}

				Skip(state, 1)
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

		p, _ := Trail(state)
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
	case func(byte) bool:
		return untilFilter(v)
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

// EOL matches the end of a line. The end of a line is one of the following:
//   a carriage return (CR, '\r')
//   a line feed (LF, '\n')
//   a carriage return + line feed (CRLF, "\r\n")
//   end of state
func EOL(state *State, result *Result) error {
	c, err := Next(state)
	if err != nil {
		result.SetToken(nil)
		return nil
	}

	if c == '\n' {
		result.SetToken([]byte{'\n'})
		state.Advance()
		return nil
	}

	if c == '\r' {
		state.Advance()
		c, err = Next(state)
		if err == nil && c == '\n' {
			state.Advance()
			result.SetToken([]byte{'\r', '\n'})
			return nil
		}
		result.SetToken([]byte{'\r'})
		return nil
	}

	return NewError("expected CR, LF, CRLF, or end of state", state.Position())
}

func calculateLineLength(state *State) (int, int) {
	i, n, cr := 0, 0, false
	for state.Request(i+1) == nil {
		c := state.Buffer()[i]
		switch {
		case c == '\n' && cr:
			return i - 1, n + 1
		case c == '\n':
			return i, n + 1
		case c == '\r':
			cr = true
			n++
		case cr:
			return i - 1, n
		}
		i++
	}
	return i, n
}

// Line matches up to a newline byte or the end of state.
func Line(state *State, result *Result) error {
	i, n := calculateLineLength(state)
	state.Request(i)
	result.SetToken(state.Buffer())
	state.Advance()
	Skip(state, n)
	return nil
}
