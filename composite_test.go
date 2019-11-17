package pars_test

import (
	"testing"

	"github.com/ktnyt/ascii"
	"github.com/ktnyt/assert"
	"github.com/ktnyt/bench"
	"github.com/ktnyt/pars"
)

func TestExact(t *testing.T) {
	s0, s1, s2 := hello, small, large
	p0, p1, p2 := []byte(s0), []byte(s1), []byte(s2)
	n0, n1, n2 := len(p0), len(p1), len(p2)
	e0 := pars.NewValueResult(s0)
	e1 := pars.NewValueResult(s1)
	e2 := pars.NewValueResult(s2)

	assert.Apply(t,
		assert.C("matching",
			assert.C(s0, MatchingCase(pars.Exact(s0), p0, e0, n0)),
			assert.C(s1, MatchingCase(pars.Exact(s1), p1, e1, n1)),
			assert.C(s2, MatchingCase(pars.Exact(s2), p2, e2, n2)),
		),
		assert.C("mismatch",
			assert.C(s0, MismatchCase(pars.Exact(p0[:5]), p0)),
			assert.C(s1, MismatchCase(pars.Exact(p1[:5]), p1)),
			assert.C(s2, MismatchCase(pars.Exact(p2[:5]), p2)),
		),
	)
}

func BenchmarkExact(b *testing.B) {
	s0, s1, s2 := hello, small, large
	p0, p1, p2 := []byte(s0), []byte(s1), []byte(s2)

	bench.Apply(b,
		bench.C("matching",
			bench.C(s0, ParserBench(pars.Exact(s0), p0)),
			bench.C(s1, ParserBench(pars.Exact(s1), p1)),
			bench.C(s2, ParserBench(pars.Exact(s2), p2)),
		),
		bench.C("mismatch",
			bench.C(s0, ParserBench(pars.Exact(p0[:5]), p0)),
			bench.C(s1, ParserBench(pars.Exact(p1[:5]), p1)),
			bench.C(s2, ParserBench(pars.Exact(p2[:5]), p2)),
		),
	)
}

func TestCount(t *testing.T) {
	p0, p1, p2 := []byte(hello), []byte(small), []byte(large)
	n := 5
	e0 := pars.NewTokenResult(p0[:n])
	e1 := pars.NewTokenResult(p1[:n])
	e2 := pars.NewTokenResult(p2[:n])
	matching := pars.Count(pars.Letter, n).Map(pars.Cat)
	mismatch := pars.Count(pars.Letter, n+1).Map(pars.Cat)

	assert.Apply(t,
		assert.C("matching",
			MatchingCase(matching, p0, e0, n),
			MatchingCase(matching, p1, e1, n),
			MatchingCase(matching, p2, e2, n),
		),
		assert.C("mismatch",
			MismatchCase(mismatch, p0),
			MismatchCase(mismatch, p1),
			MismatchCase(mismatch, p2),
		),
	)
}

func BenchmarkCount(b *testing.B) {
	p0, p1, p2 := []byte(hello), []byte(small), []byte(large)

	bench.Apply(b,
		bench.C("matching",
			ParserBench(pars.Count(pars.Letter, 5), p0),
			ParserBench(pars.Count(pars.Letter, 5), p1),
			ParserBench(pars.Count(pars.Letter, 5), p2),
		),
		bench.C("mismatch",
			ParserBench(pars.Count(pars.Letter, 6), p0),
			ParserBench(pars.Count(pars.Letter, 6), p1),
			ParserBench(pars.Count(pars.Letter, 6), p2),
		),
	)
}

func TestDelim(t *testing.T) {
	p0, p1, p2 := []byte(hello), []byte(small), []byte(large)
	l, m, n := 5, 6, 11
	e0 := pars.AsResults(p0[:l], p0[m:n])
	e1 := pars.AsResults(p1[:l], p1[m:n])
	e2 := pars.AsResults(p2[:l], p2[m:n])
	matching := pars.Delim(pars.Word(ascii.IsLetter), pars.Word(ascii.IsSpace))
	mismatch := pars.Delim(pars.Word(ascii.IsSpace), pars.Word(ascii.IsLetter))

	assert.Apply(t,
		assert.C("matching",
			MatchingCase(matching, p0, e0, n),
			MatchingCase(matching, p1, e1, n),
			MatchingCase(matching, p2, e2, n),
		),
		assert.C("mismatch",
			MismatchCase(mismatch, p0),
			MismatchCase(mismatch, p1),
			MismatchCase(mismatch, p2),
		),
	)
}

func BenchmarkDelim(b *testing.B) {
	p0, p1, p2 := []byte(hello), []byte(small), []byte(large)
	matching := pars.Delim(pars.Word(ascii.IsLetter), pars.Word(ascii.IsSpace))
	mismatch := pars.Delim(pars.Word(ascii.IsSpace), pars.Word(ascii.IsLetter))

	bench.Apply(b,
		bench.C("matching",
			ParserBench(matching, p0),
			ParserBench(matching, p1),
			ParserBench(matching, p2),
		),
		bench.C("mismatch",
			ParserBench(mismatch, p0),
			ParserBench(mismatch, p1),
			ParserBench(mismatch, p2),
		),
	)
}
