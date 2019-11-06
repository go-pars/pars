package pars_test

import (
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/pars"
)

func TestSeq(t *testing.T) {
	p0, p1 := []byte(hello), []byte(small)
	n := 5
	s := hello[:5]
	q := make([]interface{}, len(s))
	for i, c := range s {
		q[i] = c
	}
	p := pars.Seq(q...)
	e := pars.AsResults(q...)

	assert.Apply(t,
		assert.C("matching", MatchingCase(p, p0, e, n)),
		assert.C("mismatch", MismatchCase(p, p1)),
	)
}

func BenchmarkSeq(b *testing.B) {
	p0, p1 := []byte(hello), []byte(small)
	s := matchingString[:5]
	q := make([]interface{}, len(s))
	for i, c := range s {
		q[i] = c
	}
	p := pars.Seq(q...)

	b.Run("matching", benchmark(p, p0))
	b.Run("mismatch", benchmark(p, p1))
}

func TestAny(t *testing.T) {
	p0, p1, p2 := []byte(hello), []byte(small), []byte(goodbye)
	n := 5
	fst, snd := hello[:n], small[:n]
	e0, e1 := pars.NewValueResult(fst), pars.NewValueResult(snd)
	p := pars.Any(fst, snd)

	assert.Apply(t,
		assert.C("matching first", MatchingCase(p, p0, e0, n)),
		assert.C("matching second", MatchingCase(p, p1, e1, n)),
		assert.C("mismatch", MismatchCase(p, p2)),
	)
}

func BenchmarkAny(b *testing.B) {
	p0, p1, p2 := []byte(hello), []byte(small), []byte(goodbye)
	n := 5
	fst, snd := hello[:n], small[:n]
	p := pars.Any(fst, snd)

	b.Run("matching first", benchmark(p, p0))
	b.Run("matching second", benchmark(p, p1))
	b.Run("mismatch", benchmark(p, p2))
}

func TestMaybe(t *testing.T) {
	p0, p1 := []byte(hello), []byte(small)
	n := 5
	p := pars.Maybe(hello[:n])
	e := pars.NewValueResult(hello[:n])

	assert.Apply(t,
		assert.C("matching", MatchingCase(p, p0, e, n)),
		assert.C("mismatch", MatchingCase(p, p1, &pars.Result{}, 0)),
	)
}

func BenchmarkMaybe(b *testing.B) {
	p0, p1 := []byte(hello), []byte(small)
	p := pars.Maybe(hello[:5])

	b.Run("matching", benchmark(p, p0))
	b.Run("mismatch", benchmark(p, p1))
}

func TestMany(t *testing.T) {
	p := []byte(hello)
	q := make([]interface{}, len(p))
	for i, c := range p {
		q[i] = c
	}
	n0, n1 := 1, len(p)
	e0, e1 := pars.AsResults(q[:n0]...), pars.AsResults(q...)
	p0, p1, p2 := pars.Many(byte('H')), pars.Many(pars.Byte()), pars.Many('h')

	assert.Apply(t,
		assert.C("matching one", MatchingCase(p0, p, e0, n0)),
		assert.C("matching many", MatchingCase(p1, p, e1, n1)),
		assert.C("mismatch", MismatchCase(p2, p)),
	)
}

func BenchmarkMany(b *testing.B) {
	p := []byte(hello)
	p0, p1, p2 := pars.Many(byte('H')), pars.Many(pars.Byte()), pars.Many('h')
	b.Run("match one", benchmark(p0, p))
	b.Run("match many", benchmark(p1, p))
	b.Run("mismatch", benchmark(p2, p))
}

func TestCount(t *testing.T) {
	p0, p1 := []byte(goodbye), []byte(hello)
	r := []rune(goodbye)
	q := make([]interface{}, len(r))
	for i, c := range r {
		q[i] = c
	}
	e := pars.AsResults(q...)
	p := pars.Count(pars.Rune(), len(r))

	assert.Apply(t,
		assert.C("matching", MatchingCase(p, p0, e, len(p0))),
		assert.C("mismatch", MismatchCase(p, p1)),
	)
}

func BenchmarkCount(b *testing.B) {
	b.Run("matching", benchmark(pars.Count(pars.Byte(), len(matchingBytes)), matchingBytes))
	b.Run("mismatch", benchmark(pars.Count(pars.Byte(), len(mismatchBytes)), matchingBytes))
}
