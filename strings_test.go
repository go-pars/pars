package pars_test

import (
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/pars"
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
	b.Run("matching", benchmark(pars.String(hello[:5]), p0))
	b.Run("matching", benchmark(pars.String(hello[:5]), p1))
}
