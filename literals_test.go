package pars_test

import (
	"strconv"
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/pars"
)

func MatchingInt(i int) assert.F {
	s := strconv.Itoa(i)
	n := len(s)
	p := []byte(s + " is the answer")
	e := pars.NewValueResult(i)
	return MatchingCase(pars.Int, p, e, n)
}

func TestInt(t *testing.T) {
	assert.Apply(t,
		assert.C("matching",
			MatchingInt(0),
			MatchingInt(42),
			MatchingInt(-42),
		),
		assert.C("mismatch", MismatchCase(pars.Int, []byte(hello))),
	)
}

func BenchmarkInt(b *testing.B) {
	b.Run("matching", benchmark(pars.Int, []byte("42 is the answer")))
	b.Run("mismatch", benchmark(pars.Int, []byte(hello)))
}
