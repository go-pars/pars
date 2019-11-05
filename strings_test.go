package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
)

func TestString(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		if msg := matching(
			pars.String(matchingString[:5]),
			matchingBytes,
			pars.NewValueResult(matchingString[:5]),
			matchingBytes[5:],
		); msg != "" {
			t.Fatal(msg)
		}
	})
	t.Run("mismatch", func(t *testing.T) {
		if msg := mismatch(
			pars.String("matching"),
			matchingBytes,
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkString(b *testing.B) {
	b.Run("matching", benchmark(pars.String(matchingString[:5]), matchingBytes))
	b.Run("matching", benchmark(pars.String("matching"), matchingBytes))
}
