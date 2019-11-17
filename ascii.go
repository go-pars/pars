package pars

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/ktnyt/ascii"
)

// Parsers for matching ASCII character patterns.
var (
	Null    = Byte(0)
	Graphic = ByteRange(32, 126)
	Control = Any(ByteRange(0, 31), Byte(127))
	Space   = Byte(ascii.Space...)
	Upper   = ByteRange('A', 'Z')
	Lower   = ByteRange('a', 'z')
	Letter  = Any(Upper, Lower)
	Digit   = ByteRange('0', '9')
	Latin   = Any(Letter, Digit)
)

// Filter creates a Parser which will attempt to match the given ascii.Filter.
func Filter(filter ascii.Filter) Parser {
	v := reflect.ValueOf(filter)
	f := runtime.FuncForPC(v.Pointer())
	what := fmt.Sprintf("expected to match filter `%s`", f.Name())

	return func(state *State, result *Result) error {
		c, err := Next(state)
		if err != nil {
			return err
		}
		if !filter(c) {
			return NewError(what, state.Position())
		}
		state.Advance()
		result.SetToken([]byte{c})
		return nil
	}
}

// Word creates a Parser which will attempt to match a group of bytes which
// satisfy the given filter.
func Word(filter ascii.Filter) Parser {
	v := reflect.ValueOf(filter)
	f := runtime.FuncForPC(v.Pointer())
	what := fmt.Sprintf("expected to word of `%s`", f.Name())

	return func(state *State, result *Result) error {
		state.Push()
		c, err := Next(state)
		for err == nil && filter(c) {
			state.Advance()
			c, err = Next(state)
		}
		p, err := Trail(state)
		if err != nil {
			return err
		}
		if len(p) == 0 {
			return NewError(what, state.Position())
		}
		result.SetToken(p)
		return nil
	}
}
