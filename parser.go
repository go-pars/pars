package pars

import (
	"fmt"
	"io"
)

// Parser is the function signature of a parser.
type Parser func(*State, *Result) error

// Map applies the callback if the parser matches.
// Use Map to convert a result into meaningful data.
func (p Parser) Map(f Map) Parser {
	return func(state *State, result *Result) error {
		if err := p(state, result); err != nil {
			return err
		}
		f(result)
		return nil
	}
}

// Bind a specific value to a result of a parser.
func (p Parser) Bind(v interface{}) Parser {
	return func(state *State, result *Result) error {
		if err := p(state, result); err != nil {
			return err
		}
		result.Value = v
		return nil
	}
}

// ParserLike is a placeholder for types that can potentially be turned into a
// Parser with AsParser.
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
		return Bytes(p...)
	case rune:
		return Rune(p)
	case []rune:
		return Runes(p...)
	case string:
		return String(p)
	default:
		panic(fmt.Errorf("cannot make a parser from `%T`", p))
	}
}

// AsParsers attempts to create Parsers from the given ParserLike objects.
func AsParsers(q ...ParserLike) []Parser {
	p := make([]Parser, len(q))
	for i := range q {
		p[i] = AsParser(q[i])
	}
	return p
}

// Apply a parser to a State.
func Apply(p Parser, s *State) (interface{}, error) {
	r := Result{}
	if err := p(s, &r); err != nil && err != io.EOF {
		return nil, err
	}
	return r.Value, nil
}
