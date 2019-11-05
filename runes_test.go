package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
)

func TestRune(t *testing.T) {
	t.Run("no argument", func(t *testing.T) {
		if msg := matching(
			pars.Rune(),
			matchingBytes,
			pars.NewValueResult(matchingRunes[0]),
			matchingBytes[1:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("single argument", func(t *testing.T) {
		t.Run("matching", func(t *testing.T) {
			if msg := matching(
				pars.Rune(matchingRunes[0]),
				matchingBytes,
				pars.NewValueResult(matchingRunes[0]),
				matchingBytes[1:],
			); msg != "" {
				t.Fatal(msg)
			}
		})
		t.Run("mismatch", func(t *testing.T) {
			if msg := mismatch(
				pars.Rune(mismatchRunes[0]),
				matchingBytes,
			); msg != "" {
				t.Fatal(msg)
			}
		})
	})

	t.Run("many arguments", func(t *testing.T) {
		t.Run("matching first", func(t *testing.T) {
			if msg := matching(
				pars.Rune(matchingRunes[0], mismatchRunes[0]),
				matchingBytes,
				pars.NewValueResult(matchingRunes[0]),
				matchingBytes[1:],
			); msg != "" {
				t.Fatal(msg)
			}
		})
		t.Run("matching second", func(t *testing.T) {
			if msg := matching(
				pars.Rune(mismatchRunes[0], matchingRunes[0]),
				matchingBytes,
				pars.NewValueResult(matchingRunes[0]),
				matchingBytes[1:],
			); msg != "" {
				t.Fatal(msg)
			}
		})
		t.Run("mismatch", func(t *testing.T) {
			if msg := mismatch(
				pars.Rune(mismatchRunes[0], 'w'),
				matchingBytes,
			); msg != "" {
				t.Fatal(msg)
			}
		})
	})
}

func BenchmarkRune(b *testing.B) {
	b.Run("no argument", benchmark(pars.Rune(), matchingBytes))

	b.Run("single argument", combineBench(
		benchCase{"matching", benchmark(pars.Rune(matchingRunes[0]), matchingBytes)},
		benchCase{"mismatch", benchmark(pars.Rune(mismatchRunes[0]), matchingBytes)},
	))

	b.Run("many arguments", combineBench(
		benchCase{"matching first", benchmark(pars.Rune(matchingRunes[0], mismatchRunes[0]), matchingBytes)},
		benchCase{"matching second", benchmark(pars.Rune(mismatchRunes[0], matchingRunes[0]), matchingBytes)},
		benchCase{"mismatch", benchmark(pars.Rune(mismatchRunes[0], mismatchRunes[6]), matchingBytes)},
	))
}

func TestRuneRange(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		if msg := matching(
			pars.RuneRange('A', 'Z'),
			matchingBytes,
			pars.NewValueResult(matchingRunes[0]),
			matchingBytes[1:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		if msg := mismatch(
			pars.RuneRange('a', 'z'),
			matchingBytes,
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkRangeRune(b *testing.B) {
	b.Run("matching", benchmark(pars.RuneRange('A', 'Z'), matchingBytes))
	b.Run("mismatch", benchmark(pars.RuneRange('a', 'z'), matchingBytes))
}

func TestRunes(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		if msg := matching(
			pars.Runes(matchingRunes[:5]),
			matchingBytes,
			pars.NewValueResult(matchingRunes[:5]),
			matchingBytes[5:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		if msg := mismatch(
			pars.Runes([]rune("matching")),
			matchingBytes,
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkRunes(b *testing.B) {
	b.Run("matching", benchmark(pars.Runes(matchingRunes[:5]), matchingBytes))
	b.Run("mismatch", benchmark(pars.Runes(mismatchRunes[:5]), matchingBytes))
}
