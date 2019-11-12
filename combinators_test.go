package pars_test

import (
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/bench"
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

	bench.Apply(b,
		bench.C("matching", benchmark(p, p0)),
		bench.C("mismatch", benchmark(p, p1)),
	)
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

	bench.Apply(b,
		bench.C("matching first", benchmark(p, p0)),
		bench.C("matching second", benchmark(p, p1)),
		bench.C("mismatch", benchmark(p, p2)),
	)
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

	bench.Apply(b,
		bench.C("matching", benchmark(p, p0)),
		bench.C("mismatch", benchmark(p, p1)),
	)
}

func TestMany(t *testing.T) {
	p := []byte(hello)
	q := make([]interface{}, len(p))
	for i, c := range p {
		q[i] = c
	}
	n0, n1 := 1, len(p)
	e0, e1, e2 := pars.AsResults(q[:n0]...), pars.AsResults(q...), pars.AsResults()
	p0, p1, p2 := pars.Many(byte('H')), pars.Many(pars.Byte()), pars.Many('h')

	assert.Apply(t,
		assert.C("matching one", MatchingCase(p0, p, e0, n0)),
		assert.C("matching many", MatchingCase(p1, p, e1, n1)),
		assert.C("matching none", MatchingCase(p2, p, e2, 0)),
	)
}

func BenchmarkMany(b *testing.B) {
	p := []byte(hello)
	p0, p1, p2 := pars.Many(byte('H')), pars.Many(pars.Byte()), pars.Many('h')

	bench.Apply(b,
		bench.C("match one", benchmark(p0, p)),
		bench.C("match many", benchmark(p1, p)),
		bench.C("mismatch", benchmark(p2, p)),
	)
}
