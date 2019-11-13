package pars_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/bench"
	"github.com/ktnyt/pars"
)

func MatchingCase(q interface{}, p []byte, e *pars.Result, n int) assert.F {
	s := pars.FromBytes(p)
	r := &pars.Result{}
	return assert.All(
		assert.NoError(pars.AsParser(q)(s, r)),
		assert.Equal(r.Token, e.Token),
		assert.Equal(r.Value, e.Value),
		assert.Equal(r.Children, e.Children),
		assert.Equal(s.Dump(), p[n:]),
	)
}

func MismatchCase(q interface{}, p []byte) assert.F {
	s := pars.FromBytes(p)
	r := &pars.Result{}
	e := &pars.Result{}
	return assert.All(
		assert.IsError(pars.AsParser(q)(s, r)),
		assert.Equal(r.Token, e.Token),
		assert.Equal(r.Value, e.Value),
		assert.Equal(r.Children, e.Children),
		assert.Equal(s.Dump(), p),
	)
}

func noerror(err error) string {
	if err != nil {
		return fmt.Sprintf("unexpected error: %s", err.Error())
	}
	return ""
}

func iserror(err error) string {
	if err == nil {
		return "expected error"
	}
	return ""
}

func equals(act, exp interface{}) string {
	if !reflect.DeepEqual(exp, act) {
		return fmt.Sprintf("expected: %#v\n  actual: %#v", exp, act)
	}
	return ""
}

func try(msgs ...string) string {
	for _, msg := range msgs {
		if msg != "" {
			return msg
		}
	}
	return ""
}

var hello = "Hello world!"
var small = "Small world!"
var large = "Large world!"

var matchingString = "Hello world!"
var matchingBytes = []byte(matchingString)
var matchingRunes = []rune(matchingString)

var mismatchString = "Goodbye world!"
var mismatchBytes = []byte(mismatchString)
var mismatchRunes = []rune(mismatchString)

func matching(p pars.Parser, in []byte, e *pars.Result, out []byte) string {
	s := pars.FromBytes(in)
	r := &pars.Result{}
	return try(
		noerror(p(s, r)),
		equals(r.Token, e.Token),
		equals(r.Value, e.Value),
		equals(r.Children, e.Children),
		equals(s.Dump(), out),
	)
}

func mismatch(p pars.Parser, in []byte) string {
	s := pars.FromBytes(in)
	r := &pars.Result{}
	return try(
		iserror(p(s, r)),
		equals(r, &pars.Result{}),
		equals(s.Dump(), in),
	)
}

func benchmark(p pars.Parser, in []byte) bench.F {
	s := pars.FromBytes(in)
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	}
}
