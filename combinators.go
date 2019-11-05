package pars

import (
	"fmt"
	"reflect"
	"strings"
)

func typeRep(q interface{}) string {
	return reflect.TypeOf(q).String()
}

func typeReps(qs []interface{}) []string {
	r := make([]string, len(qs))
	for i, q := range qs {
		r[i] = typeRep(q)
	}
	return r
}

func Seq(qs ...interface{}) Parser {
	ps := AsParsers(qs...)
	name := fmt.Sprintf("Seq(%s)", strings.Join(typeReps(qs), ", "))

	return func(state *State, result *Result) error {
		v := make([]Result, len(ps))

		state.Push()
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

func Any(qs ...interface{}) Parser {
	ps := AsParsers(qs...)
	name := fmt.Sprintf("Any(%s)", strings.Join(typeReps(qs), ", "))

	return func(state *State, result *Result) error {
		state.Push()
		for _, p := range ps {
			if p(state, result) == nil {
				state.Drop()
				return nil
			}
		}
		state.Pop()

		return NewParserError(name, state.Position())
	}
}

func Maybe(q interface{}) Parser {
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

func Many(q interface{}) Parser {
	p := AsParser(q)
	name := fmt.Sprintf("Many(%s)", typeRep(q))

	return func(state *State, result *Result) error {
		v := make([]Result, 1, 5)

		state.Push()
		if err := p(state, &v[0]); err != nil {
			state.Pop()
			return NewTraceError(name, err)
		}
		state.Drop()

		state.Push()
		for p(state, result) == nil {
			state.Drop()
			state.Push()
			v = append(v, *result)
			*result = Result{}
		}
		state.Pop()

		result.SetChildren(v)
		return nil
	}
}

func Count(q interface{}, n int) Parser {
	p := AsParser(q)
	name := fmt.Sprintf("Count(%s, %d)", typeRep(q), n)

	return func(state *State, result *Result) error {
		v := make([]Result, n)

		state.Push()
		for i := 0; i < n; i++ {
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
