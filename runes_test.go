package pars_test

import (
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/pars"
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

	b.Run("no argument", benchmark(pars.Rune(), p0))

	b.Run("single argument", combineBench(
		benchCase{"matching", benchmark(pars.Rune(r0[0]), p0)},
		benchCase{"mismatch", benchmark(pars.Rune(r0[0]), p1)},
	))

	b.Run("many arguments", combineBench(
		benchCase{"matching first", benchmark(pars.Rune(r0[0], r1[0]), p0)},
		benchCase{"matching second", benchmark(pars.Rune(r1[0], r0[0]), p0)},
		benchCase{"mismatch", benchmark(pars.Rune(r0[0]), p1)},
	))
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
	b.Run("matching", benchmark(pars.RuneRange('A', 'Z'), p))
	b.Run("mismatch", benchmark(pars.RuneRange('a', 'z'), p))
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
	b.Run("matching", benchmark(pars.Runes([]rune(hello)[:5]), p0))
	b.Run("mismatch", benchmark(pars.Runes([]rune(hello)[:5]), p1))
}
