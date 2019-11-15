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
		if err := p(state, result); err != nil {
			return err
		}
		return f(result)
	}
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
