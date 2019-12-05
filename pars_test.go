package pars_test

import (
	"testing"

	"gopkg.in/ktnyt/assert.v1"
	"gopkg.in/ktnyt/bench.v1"
	"gopkg.in/ktnyt/pars.v2"
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

func ParserBench(p pars.Parser, in []byte) bench.F {
	s := pars.FromBytes(in)
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	}
}

var hello = "Hello world!"
var small = "Small world!"
var large = "Large world!"
