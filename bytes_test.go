package pars_test

import (
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/bench"
	"github.com/ktnyt/pars"
)

func TestByte(t *testing.T) {
	p := []byte(hello)
	n := 1
	e := pars.NewTokenResult(p[:n])

	assert.Apply(t,
		assert.C("no argument", MatchingCase(pars.Byte(), p, e, n)),
		assert.C("single argument",
			assert.C("matching", MatchingCase(pars.Byte('H'), p, e, n)),
			assert.C("mismatch", MismatchCase(pars.Byte('h'), p)),
		),
		assert.C("multiple arguments",
			assert.C("match first", MatchingCase(pars.Byte('H', 'h'), p, e, n)),
			assert.C("match second", MatchingCase(pars.Byte('h', 'H'), p, e, n)),
			assert.C("mismatch", MismatchCase(pars.Byte('h', 'w'), p)),
		),
	)
}

func BenchmarkByte(b *testing.B) {
	p0, p1 := []byte(hello), []byte(small)

	bench.Apply(b,
		bench.C("no argument", ParserBench(pars.Byte(), p0)),
		bench.C("single argument",
			bench.C("matching", ParserBench(pars.Byte(p0[0]), p0)),
			bench.C("mismatch", ParserBench(pars.Byte(p0[0]), p1)),
		),
		bench.C("many arguments",
			bench.C("matching first", ParserBench(pars.Byte(p0[0], p1[0]), p0)),
			bench.C("matching second", ParserBench(pars.Byte(p1[0], p0[0]), p0)),
			bench.C("mismatch", ParserBench(pars.Byte(p0[0]), p1)),
		),
	)

}

func TestByteRange(t *testing.T) {
	p := []byte(hello)
	n := 1
	e := pars.NewTokenResult(p[:n])

	assert.Apply(t,
		assert.C("matching", MatchingCase(pars.ByteRange('A', 'Z'), p, e, n)),
		assert.C("mismatch", MismatchCase(pars.ByteRange('a', 'z'), p)),
	)
}

func BenchmarkRangeByte(b *testing.B) {
	p := []byte(hello)
	bench.Apply(b,
		bench.C("matching", ParserBench(pars.ByteRange('A', 'Z'), p)),
		bench.C("mismatch", ParserBench(pars.ByteRange('a', 'z'), p)),
	)
}

func TestBytes(t *testing.T) {
	p := []byte(hello)
	n := 5
	e := pars.NewTokenResult(p[:n])

	assert.Apply(t,
		assert.C("matching", MatchingCase(pars.Bytes(p[:n]), p, e, n)),
		assert.C("mismatch", MismatchCase(pars.Bytes(p[n:]), p)),
	)
}

func BenchmarkBytes(b *testing.B) {
	p0, p1 := []byte(hello), []byte(small)
	bench.Apply(b,
		bench.C("matching", ParserBench(pars.Bytes(p0[:5]), p0)),
		bench.C("mismatch", ParserBench(pars.Bytes(p0[:5]), p1)),
	)
}
