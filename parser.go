package pars

import "fmt"

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

// ParserLike is a placeholder for types that can potentially be converted to a
// Parser with the AsParser function.
type ParserLike interface{}

// AsParser attempts to create a Parser from a ParserLike object.
func AsParser(q ParserLike) Parser {
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
	default:
		panic(fmt.Errorf("cannot convert type `%T` to a parser", p))
	}
}

// AsParsers applies the AsParser function to each argument.
func AsParsers(qs ...ParserLike) []Parser {
	ps := make([]Parser, len(qs))
	for i, q := range qs {
		ps[i] = AsParser(q)
	}
	return ps
}
