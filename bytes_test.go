package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
)

func TestByte(t *testing.T) {
	t.Run("no argument", func(t *testing.T) {
		if msg := matching(
			pars.Byte(),
			matchingBytes,
			pars.NewTokenResult(matchingBytes[:1]),
			matchingBytes[1:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("single argument", func(t *testing.T) {
		t.Run("matching", func(t *testing.T) {
			if msg := matching(
				pars.Byte(matchingBytes[0]),
				matchingBytes,
				pars.NewTokenResult(matchingBytes[:1]),
				matchingBytes[1:],
			); msg != "" {
				t.Fatal(msg)
			}
		})

		t.Run("mismatch", func(t *testing.T) {
			if msg := mismatch(
				pars.Byte(mismatchBytes[0]),
				matchingBytes,
			); msg != "" {
				t.Fatal(msg)
			}
		})
	})

	t.Run("multiple argument", func(t *testing.T) {
		t.Run("matching first", func(t *testing.T) {
			if msg := matching(
				pars.Byte(matchingBytes[0], mismatchBytes[0]),
				matchingBytes,
				pars.NewTokenResult(matchingBytes[:1]),
				matchingBytes[1:],
			); msg != "" {
				t.Fatal(msg)
			}
		})

		t.Run("matching second", func(t *testing.T) {
			if msg := matching(
				pars.Byte(mismatchBytes[0], matchingBytes[0]),
				matchingBytes,
				pars.NewTokenResult(matchingBytes[:1]),
				matchingBytes[1:],
			); msg != "" {
				t.Fatal(msg)
			}
		})

		t.Run("mismatch", func(t *testing.T) {
			if msg := mismatch(
				pars.Byte(mismatchBytes[0], mismatchBytes[6]),
				matchingBytes,
			); msg != "" {
				t.Fatal(msg)
			}
		})
	})
}

func BenchmarkByte(b *testing.B) {
	b.Run("no argument", benchmark(pars.Byte(), matchingBytes))

	b.Run("single argument", combineBench(
		benchCase{"matching", benchmark(pars.Byte(matchingBytes[0]), matchingBytes)},
		benchCase{"mismatch", benchmark(pars.Byte(mismatchBytes[0]), matchingBytes)},
	))

	b.Run("many arguments", combineBench(
		benchCase{"matching first", benchmark(pars.Byte(matchingBytes[0], mismatchBytes[0]), matchingBytes)},
		benchCase{"matching second", benchmark(pars.Byte(mismatchBytes[0], matchingBytes[0]), matchingBytes)},
		benchCase{"mismatch", benchmark(pars.Byte(mismatchBytes[0]), matchingBytes)},
	))
}

func TestByteRange(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		if msg := matching(
			pars.ByteRange('A', 'Z'),
			matchingBytes,
			pars.NewTokenResult(matchingBytes[:1]),
			matchingBytes[1:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		if msg := mismatch(
			pars.ByteRange('a', 'z'),
			matchingBytes,
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkRangeByte(b *testing.B) {
	b.Run("matching", benchmark(pars.ByteRange('A', 'Z'), matchingBytes))
	b.Run("mismatch", benchmark(pars.ByteRange('a', 'z'), matchingBytes))
}

func TestBytes(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		if msg := matching(
			pars.Bytes(matchingBytes[:5]),
			matchingBytes,
			pars.NewTokenResult(matchingBytes[:5]),
			matchingBytes[5:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		if msg := mismatch(
			pars.Bytes([]byte("matching")),
			matchingBytes,
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkBytes(b *testing.B) {
	b.Run("matching", benchmark(pars.Bytes(matchingBytes[:5]), matchingBytes))
	b.Run("mismatch", benchmark(pars.Bytes([]byte("matching")), matchingBytes))
}
