package pars_test

import (
	"testing"

	"gopkg.in/ktnyt/assert.v1"
	"gopkg.in/ktnyt/bench.v1"
	"gopkg.in/ktnyt/pars.v2"
)

func TestRune(t *testing.T) {
	p := []byte(hello)
	n := 1
	e := pars.NewValueResult([]rune(hello)[0])

	assert.Apply(t,
		assert.C("no argument", MatchingCase(pars.Rune(), p, e, n)),
		assert.C("single argument",
			assert.C("matching", MatchingCase(pars.Rune('H'), p, e, n)),
			assert.C("mismatch", MismatchCase(pars.Rune('h'), p)),
		),
		assert.C("multiple arguments",
			assert.C("match first", MatchingCase(pars.Rune('H', 'h'), p, e, n)),
			assert.C("match second", MatchingCase(pars.Rune('h', 'H'), p, e, n)),
			assert.C("mismatch", MismatchCase(pars.Rune('h', 'w'), p)),
		),
	)
}

func BenchmarkRune(b *testing.B) {
	p0, p1 := []byte(hello), []byte(small)
	r0, r1 := []rune(hello), []rune(small)

	bench.Apply(b,
		bench.C("no argument", ParserBench(pars.Rune(), p0)),
		bench.C("single argument",
			bench.C("matching", ParserBench(pars.Rune(r0[0]), p0)),
			bench.C("mismatch", ParserBench(pars.Rune(r0[0]), p1)),
		),
		bench.C("many arguments",
			bench.C("matching first", ParserBench(pars.Rune(r0[0], r1[0]), p0)),
			bench.C("matching second", ParserBench(pars.Rune(r1[0], r0[0]), p0)),
			bench.C("mismatch", ParserBench(pars.Rune(r0[0]), p1)),
		),
	)
}

func TestRuneRange(t *testing.T) {
	p := []byte(hello)
	n := 1
	e := pars.NewValueResult([]rune(hello)[0])

	assert.Apply(t,
		assert.C("matching", MatchingCase(pars.RuneRange('A', 'Z'), p, e, n)),
		assert.C("mismatch", MismatchCase(pars.RuneRange('a', 'z'), p)),
	)
}

func BenchmarkRangeRune(b *testing.B) {
	p := []byte(hello)

	bench.Apply(b,
		bench.C("matching", ParserBench(pars.RuneRange('A', 'Z'), p)),
		bench.C("mismatch", ParserBench(pars.RuneRange('a', 'z'), p)),
	)
}

func TestRunes(t *testing.T) {
	p0, p1 := []byte(hello), []byte(small)
	r := []rune(hello)
	n := 5
	e := pars.NewValueResult([]rune(hello)[:5])

	assert.Apply(t,
		assert.C("matching", MatchingCase(pars.Runes(r[:n]), p0, e, n)),
		assert.C("mismatch", MismatchCase(pars.Runes(r[:n]), p1)),
	)
}

func BenchmarkRunes(b *testing.B) {
	p0, p1 := []byte(hello), []byte(small)

	bench.Apply(b,
		bench.C("matching", ParserBench(pars.Runes([]rune(hello)[:5]), p0)),
		bench.C("mismatch", ParserBench(pars.Runes([]rune(hello)[:5]), p1)),
	)
}
