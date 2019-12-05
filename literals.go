package pars

import (
	"fmt"
	"strconv"

	"gopkg.in/ktnyt/ascii.v1"
)

func convertInt(state *State, result *Result) error {
	p, err := Trail(state)
	if err != nil {
		return err
	}
	n, err := strconv.Atoi(string(p))
	if err != nil {
		return err
	}
	result.SetValue(n)
	return nil
}

func convertNumber(state *State, result *Result) error {
	p, err := Trail(state)
	if err != nil {
		return err
	}
	n, err := strconv.ParseFloat(string(p), 64)
	if err != nil {
		return err
	}
	result.SetValue(n)
	return nil
}

// Int will match an integer string and return its numerical representation.
// An integer is defined to be as follows in EBNF:
//
//   zero     = `0`
//   non-zero = `1` | `2` | `3` | `4` | `5` | `6` | `7` | `8` | `9`
//   digit    = zero | non-zero
//   integer  = [-], zero | digit, { digit }
//
// This implementation is optimized so the parser will first scan as far as
// possible to match a valid integer and then retrieve a block of bytes and
// convert it to an `int` via strconv.Atoi.
func Int(state *State, result *Result) error {
	// Scan forwards from the current position.
	state.Push()

	c, err := Next(state)
	if err != nil {
		return NewNestedError("Int", err)
	}

	// Test the first byte for a possible negative sign.
	if c == '-' {
		state.Advance()
		c, err = Next(state)
		if err != nil {
			return NewNestedError("Int", err)
		}
	}

	// The byte is not a digit so return an error.
	if !ascii.IsDigit(c) {
		state.Pop()
		return NewError("expected an integer", state.Position())
	}

	// The byte is a `0` so immediately return a 0.
	if c == '0' {
		state.Advance()
		state.Drop()
		result.SetValue(0)
		return nil
	}

	// Continually match digits.
	for err == nil && ascii.IsDigit(c) {
		state.Advance()
		c, err = Next(state)
	}

	if err := convertInt(state, result); err != nil {
		return NewNestedError("Int", err)
	}

	return nil
}

// Number will match a floating point number.
// A number is defined to be as follows in EBNF:
//
//   zero     = `0`
//   non-zero = `1` | `2` | `3` | `4` | `5` | `6` | `7` | `8` | `9`
//   digit    = zero | non-zero
//   integer  = [ `-` ], zero | digit, { digit }
//   fraction = `.`, digit, { digit }
//   exponent = ( `e` | `E` ), [ ( `-` | `+` ) ], integer
//   number   = [ `-` ],  integer, [ fraction ], [ exponent ]
//
// This implementation is optimized so the parser will first scan as far as
// possible to match a valid number and then retrieve a block of bytes and
// convert it to a `float64` via strconv.ParseFloat.
func Number(state *State, result *Result) error {
	// Scan forwards from the current position.
	state.Push()

	c, err := Next(state)
	if err != nil {
		return NewNestedError("Number", err)
	}

	// Test the first byte for a possible negative sign.
	if c == '-' {
		state.Advance()
		c, err = Next(state)
		if err != nil {
			state.Pop()
			return NewNestedError("Number", err)
		}
	}

	// Test the byte for a digit.
	if !ascii.IsDigit(c) {
		state.Pop()
		return NewError("expected a number", state.Position())
	}

	// Process the integer part.
	// Advance more than once if the first digit is not zero.
	if c == '0' {
		state.Advance()
		c, err = Next(state)
	} else {
		state.Advance()
		c, err = Next(state)

		// Continually match digits.
		for err == nil && ascii.IsDigit(c) {
			state.Advance()
			c, err = Next(state)
		}
	}

	// Process the fraction part.
	if err == nil && c == '.' {
		// The parser may need to backtrack to this position.
		state.Push()

		state.Advance()
		c, err = Next(state)
		if err != nil {
			state.Pop()
			if err := convertNumber(state, result); err != nil {
				return NewNestedError("Number", err)
			}
			return nil
		}

		if !ascii.IsDigit(c) {
			state.Pop()
			return NewError("expected a number", state.Position())
		}

		state.Advance()
		c, err = Next(state)

		// Continually match digits.
		for err == nil && ascii.IsDigit(c) {
			state.Advance()
			c, err = Next(state)
		}

		// Reached the full extent of the fraction.
		state.Drop()
	}

	// Process the exponent part.
	if err == nil && (c == 'e' || c == 'E') {
		// The parser may need to backtrack to this position.
		state.Push()

		state.Advance()
		c, err = Next(state)
		if err != nil {
			state.Pop()
			if err := convertNumber(state, result); err != nil {
				return NewNestedError("Number", err)
			}
			return nil
		}

		// Test the byte for a possible positive or negative sign.
		if c == '-' || c == '+' {
			state.Advance()
			c, err = Next(state)
			if err != nil {
				state.Pop()
				if err := convertNumber(state, result); err != nil {
					return NewNestedError("Number", err)
				}
				return nil
			}
		}

		// There are no digits so backtrack and return.
		if !ascii.IsDigit(c) {
			state.Pop()
			if err := convertNumber(state, result); err != nil {
				return NewNestedError("Number", err)
			}
			return nil
		}

		state.Advance()
		c, err = Next(state)

		// Continually match digits.
		for err == nil && ascii.IsDigit(c) {
			state.Advance()
			c, err = Next(state)
		}

		// Reached the full extent of the exponent.
		state.Drop()
	}

	if err := convertNumber(state, result); err != nil {
		return NewNestedError("Number", err)
	}
	return nil
}

// Between creates a Parser which will attempt to match a sequence of bytes
// between the given bytes. If a backslash appears in the middle of the string,
// the byte immediately following will be skipped.
func Between(l, r byte) Parser {
	name := fmt.Sprintf("Between(%s, %s)", ascii.Rep(l), ascii.Rep(r))
	whatL := fmt.Sprintf("expected opening `%c`", l)
	whatR := fmt.Sprintf("expected closing `%c`", r)

	return func(state *State, result *Result) error {
		state.Push()

		c, err := Next(state)
		if err != nil {
			return NewNestedError(name, err)
		}
		if c != l {
			return NewError(whatL, state.Position())
		}

		state.Advance()

		c, err = Next(state)
		for err == nil && c != r {
			if c == '\\' {
				state.Advance()
				_, err = Next(state)
				if err != nil {
					state.Pop()
					return NewError(whatR, state.Position())
				}
			}
			state.Advance()
			c, err = Next(state)
		}

		if err != nil {
			state.Pop()
			return NewError(whatR, state.Position())
		}

		p, err := Trail(state)
		if err != nil {
			return NewNestedError(name, err)
		}
		if err := Skip(state, 1); err != nil {
			return NewNestedError(name, err)
		}
		result.SetToken(p[1:])
		return nil
	}
}

// Quoted creates a Parser which will attempt to match a sequence of bytes
// flanked by the given byte. This Parser is equivalent to the following:
//   Between(c, c)
func Quoted(c byte) Parser { return Between(c, c) }
