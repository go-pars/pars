package pars

import (
	"fmt"

	"github.com/ktnyt/ascii"
)

// Parser is the function signature of a parser.
type Parser func(*State, *Result) error

// Map is the function signature for a result mapper.
type Map func(result *Result) error

// Map applies the callback if the parser matches.
func (p Parser) Map(f Map) Parser {
	return func(state *State, result *Result) error {
		state.Push()
		if err := p(state, result); err != nil {
			state.Pop()
			return err
		}
		state.Drop()
		return f(result)
	}
}

// Child will map to the i'th child of the result.
func (p Parser) Child(i int) Parser { return p.Map(Child(i)) }

// Children will keep the children associated to the given indices.
func (p Parser) Children(indices ...int) Parser {
	return p.Map(Children(indices...))
}

// ToString will convert the Token field to a string Value.
func (p Parser) ToString() Parser { return p.Map(ToString) }

// Bind will bind the given value as the parser result value.
func (p Parser) Bind(v interface{}) Parser {
	return func(state *State, result *Result) error {
		if err := p(state, result); err != nil {
			return err
		}
		result.SetValue(v)
		return nil
	}
}

// Error will modify the Parser to return the given error if the Parser returns
// an error.
func (p Parser) Error(alt error) Parser {
	return func(state *State, result *Result) error {
		if err := p(state, result); err != nil {
			return alt
		}
		return nil
	}
}

// Parse the given state using the parser and return the Result.
func (p Parser) Parse(s *State) (Result, error) {
	r := Result{}
	err := p(s, &r)
	return r, err
}

// AsParser attempts to create a Parser for a given argument.
func AsParser(q interface{}) Parser {
	switch p := q.(type) {
	case Parser:
		return p
	case func(*State, *Result) error:
		return p
	case *Parser:
		return func(state *State, result *Result) error {
			return (*p)(state, result)
		}
	case byte:
		return Byte(p)
	case []byte:
		return Bytes(p)
	case rune:
		return Rune(p)
	case []rune:
		return Runes(p)
	case string:
		return String(p)
	case ascii.Filter:
		return Filter(p)
	default:
		panic(fmt.Errorf("cannot convert type `%T` to a parser", p))
	}
}

// AsParsers applies the AsParser function to each argument.
func AsParsers(qs ...interface{}) []Parser {
	ps := make([]Parser, len(qs))
	for i, q := range qs {
		ps[i] = AsParser(q)
	}
	return ps
}
