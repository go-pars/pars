package pars_test

import (
	"testing"

	"gopkg.in/ktnyt/ascii.v1"
	"gopkg.in/ktnyt/assert.v1"
	"gopkg.in/ktnyt/bench.v1"
	"gopkg.in/ktnyt/pars.v2"
)

func TestUntil(t *testing.T) {
	n := 5
	p0, p1, p2 := []byte(hello), []byte(small), []byte(large)
	e0 := pars.NewTokenResult(p0[:n])
	e1 := pars.NewTokenResult(p1[:n])
	e2 := pars.NewTokenResult(p2[:n])

	assert.Apply(t,
		assert.C("matching space byte",
			assert.C(hello, MatchingCase(pars.Until(' '), p0, e0, n)),
			assert.C(small, MatchingCase(pars.Until(' '), p1, e1, n)),
			assert.C(large, MatchingCase(pars.Until(' '), p2, e2, n)),
		),
		assert.C("matching space filter",
			assert.C(hello, MatchingCase(pars.Until(ascii.IsSpaceFilter), p0, e0, n)),
			assert.C(small, MatchingCase(pars.Until(ascii.IsSpaceFilter), p1, e1, n)),
			assert.C(large, MatchingCase(pars.Until(ascii.IsSpaceFilter), p2, e2, n)),
		),
		assert.C("mismatch",
			MismatchCase(pars.Until('\n'), p0),
			MismatchCase(pars.Until('\n'), p1),
			MismatchCase(pars.Until('\n'), p2),
		),
	)
}

func BenchmarkUntil(b *testing.B) {
	p0, p1, p2 := []byte(hello), []byte(small), []byte(large)

	bench.Apply(b,
		bench.C("matching space byte",
			bench.C(hello, ParserBench(pars.Until(' '), p0)),
			bench.C(small, ParserBench(pars.Until(' '), p1)),
			bench.C(large, ParserBench(pars.Until(' '), p2)),
		),
		bench.C("matching space filter",
			bench.C(hello, ParserBench(pars.Until(ascii.IsSpaceFilter), p0)),
			bench.C(small, ParserBench(pars.Until(ascii.IsSpaceFilter), p1)),
			bench.C(large, ParserBench(pars.Until(ascii.IsSpaceFilter), p2)),
		),
		bench.C("matching",
			ParserBench(pars.Until(' '), p0),
			ParserBench(pars.Until(' '), p1),
			ParserBench(pars.Until(' '), p2),
		),
		bench.C("mismatch",
			ParserBench(pars.Until('\n'), p0),
			ParserBench(pars.Until('\n'), p1),
			ParserBench(pars.Until('\n'), p2),
		),
	)
}
