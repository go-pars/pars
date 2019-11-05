package pars_test

import (
	"strconv"
	"testing"

	"github.com/ktnyt/pars"
)

func TestInt(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		ns := []int{0, 42, -42}
		for _, n := range ns {
			e := strconv.Itoa(n)
			t.Run(e, func(t *testing.T) {
				p := []byte(e + " is the answer")
				if msg := matching(
					pars.Int,
					p,
					pars.NewValueResult(n),
					p[len(e):],
				); msg != "" {
					t.Fatal(msg)
				}
			})
		}
	})
}

func BenchmarkInt(b *testing.B) {
	b.Run("matching", benchmark(pars.Int, []byte("42 is the answer")))
	b.Run("mismatch", benchmark(pars.Int, mismatchBytes))
}
