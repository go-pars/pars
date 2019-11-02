package pars

import (
	"fmt"
	"reflect"
	"strings"
)

func typeRep(q ParserLike) string {
	return reflect.TypeOf(q).String()
}

func typeReps(qs []ParserLike) []string {
	r := make([]string, len(qs))
	for i, q := range qs {
		r[i] = typeRep(q)
	}
	return r
}

func Seq(qs ...ParserLike) Parser {
	ps := AsParsers(qs...)
	name := fmt.Sprintf("Seq(%s)", strings.Join(typeReps(qs), ", "))

	return func(state *State, result *Result) error {
		state.Push()
		v := make([]Result, len(ps))
		for i, p := range ps {
			if err := p(state, &v[i]); err != nil {
				state.Pop()
				return NewTraceError(name, err)
			}
		}
		state.Drop()
		result.SetChildren(v)
		return nil
	}
}

func Any(qs ...ParserLike) Parser {
	ps := AsParsers(qs...)
	name := fmt.Sprintf("Any(%s)", strings.Join(typeReps(qs), ", "))

	return func(state *State, result *Result) error {
		state.Push()
		for i, p := range ps {
			if errs[i] := p(state, result); errs[i] == nil {
				state.Drop()
				return nil
			}
		}
		state.Pop()
		return NewParserError(name, state.Position())
	}
}

func Maybe(q ParserLike) Parser {
	p := AsParser(q)

	return func(state *State, result *Result) error {
  	state.Push()
  	if p(state, result) != nil {
    	state.Pop()
    	return nil
  	}
  	state.Drop()
		return nil
	}
}
