package pars_test

import (
	"testing"

	"gopkg.in/ktnyt/assert.v1"
	"gopkg.in/ktnyt/pars.v2"
)

func TestString(t *testing.T) {
	p0, p1 := []byte(hello), []byte(small)
	n := 5
	e := pars.NewValueResult(hello[:n])

	assert.Apply(t,
		assert.C("matching", MatchingCase(pars.String(hello[:n]), p0, e, n)),
		assert.C("mismatch", MismatchCase(pars.String(hello[:n]), p1)),
	)
}

func BenchmarkString(b *testing.B) {
	p0, p1 := []byte(hello), []byte(small)
	b.Run("matching", ParserBench(pars.String(hello[:5]), p0))
	b.Run("matching", ParserBench(pars.String(hello[:5]), p1))
}
