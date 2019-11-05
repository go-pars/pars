package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
)

func TestSeq(t *testing.T) {
	s := matchingString[:5]
	q := make([]interface{}, len(s))
	for i, c := range s {
		q[i] = c
	}

	t.Run("matching", func(t *testing.T) {
		if msg := matching(
			pars.Seq(q...),
			matchingBytes,
			pars.AsResults(q...),
			matchingBytes[5:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		if msg := mismatch(
			pars.Seq(q...),
			mismatchBytes,
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkSeq(b *testing.B) {
	s := matchingString[:5]
	q := make([]interface{}, len(s))
	for i, c := range s {
		q[i] = c
	}
	p := pars.Seq(q...)

	b.Run("matching", benchmark(p, matchingBytes))
	b.Run("mismatch", benchmark(p, mismatchBytes))
}

func TestAny(t *testing.T) {
	first := matchingString[:5]
	second := mismatchString[:7]

	t.Run("matching first", func(t *testing.T) {
		if msg := matching(
			pars.Any(first, second),
			matchingBytes,
			pars.NewValueResult(first),
			matchingBytes[5:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("matching second", func(t *testing.T) {
		if msg := matching(
			pars.Any(second, first),
			matchingBytes,
			pars.NewValueResult(first),
			matchingBytes[5:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		if msg := mismatch(
			pars.Any(first, second),
			[]byte("hello world!"),
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkAny(b *testing.B) {
	first := matchingBytes[:5]
	second := mismatchBytes[:7]

	b.Run("matching first", benchmark(pars.Any(first, second), matchingBytes))
	b.Run("matching second", benchmark(pars.Any(second, first), matchingBytes))
	b.Run("mismatch", benchmark(pars.Any(first, second), []byte("hello world!")))
}

func TestMaybe(t *testing.T) {
	parser := pars.Maybe(matchingString[:5])

	t.Run("matching", func(t *testing.T) {
		if msg := matching(
			parser,
			matchingBytes,
			pars.NewValueResult(matchingString[:5]),
			matchingBytes[5:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		if msg := matching(
			parser,
			mismatchBytes,
			&pars.Result{},
			mismatchBytes,
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkMaybe(b *testing.B) {
	parser := pars.Maybe(matchingString[:5])

	b.Run("matching", benchmark(parser, matchingBytes))
	b.Run("mismatch", benchmark(parser, mismatchBytes))
}

func TestMany(t *testing.T) {
	q := make([]interface{}, len(matchingBytes))
	for i, c := range matchingBytes {
		q[i] = c
	}

	t.Run("matching one", func(t *testing.T) {
		if msg := matching(
			pars.Many(pars.ByteRange('A', 'Z')),
			matchingBytes,
			pars.AsResults(q[:1]...),
			matchingBytes[1:],
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("matching many", func(t *testing.T) {
		if msg := matching(
			pars.Many(pars.Byte()),
			matchingBytes,
			pars.AsResults(q...),
			[]byte{},
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		if msg := mismatch(
			pars.Many(pars.ByteRange('a', 'z')),
			matchingBytes,
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkMany(b *testing.B) {
	b.Run("match one", benchmark(pars.Many(pars.ByteRange('A', 'Z')), matchingBytes))
	b.Run("match many", benchmark(pars.Many(pars.Byte()), matchingBytes))
	b.Run("mismatch", benchmark(pars.Many(pars.ByteRange('a', 'z')), matchingBytes))
}

func TestCount(t *testing.T) {
	q := make([]interface{}, len(matchingString))
	for i, c := range matchingString {
		q[i] = c
	}

	t.Run("matching", func(t *testing.T) {
		if msg := matching(
			pars.Count(pars.Rune(), len(matchingBytes)),
			matchingBytes,
			pars.AsResults(q...),
			[]byte{},
		); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		if msg := mismatch(
			pars.Count(pars.Rune(), len(matchingBytes)),
			matchingBytes[:5],
		); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkCount(b *testing.B) {
	b.Run("matching", benchmark(pars.Count(pars.Byte(), len(matchingBytes)), matchingBytes))
	b.Run("mismatch", benchmark(pars.Count(pars.Byte(), len(mismatchBytes)), matchingBytes))
}
