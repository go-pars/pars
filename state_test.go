package pars_test

import (
	"io"
	"testing"

	"github.com/ktnyt/pars"
)

func TestState(t *testing.T) {
	t.Run("Read", func(t *testing.T) {
		s := pars.FromBytes(matchingBytes)
		p := make([]byte, len(matchingBytes))
		n, err := s.Read(p)

		if msg := equals(n, len(matchingBytes)); msg != "" {
			t.Fatal(msg)
		}
		if msg := noerror(err); msg != "" {
			t.Fatal(msg)
		}
		if msg := equals(p, matchingBytes); msg != "" {
			t.Fatal(msg)
		}
		if msg := equals(s.Dump(), []byte{}); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("Want", func(t *testing.T) {
		s := pars.FromBytes(matchingBytes)

		if msg := noerror(s.Want(len(matchingBytes))); msg != "" {
			t.Fatal(msg)
		}
		if msg := equals(s.Want(len(matchingBytes)+1), io.EOF); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("Advance", func(t *testing.T) {
		s := pars.FromBytes(matchingBytes)

		if msg := noerror(s.Want(1)); msg != "" {
			t.Fatal(msg)
		}
		s.Advance()
		if msg := equals(s.Dump(), matchingBytes[1:]); msg != "" {
			t.Fatal(msg)
		}

		if msg := noerror(s.Want(5)); msg != "" {
			t.Fatal(msg)
		}
		s.Advance()
		if msg := equals(s.Dump(), matchingBytes[6:]); msg != "" {
			t.Fatal(msg)
		}
	})

	t.Run("Stack", func(t *testing.T) {
		s := pars.FromBytes(matchingBytes)

		s.Push()
		if msg := noerror(s.Want(1)); msg != "" {
			t.Fatal(msg)
		}
		s.Advance()
		if msg := equals(s.Dump(), matchingBytes[1:]); msg != "" {
			t.Fatal(msg)
		}
		s.Push()
		if msg := noerror(s.Want(5)); msg != "" {
			t.Fatal(msg)
		}
		s.Advance()
		if msg := equals(s.Dump(), matchingBytes[6:]); msg != "" {
			t.Fatal(msg)
		}
		s.Pop()
		if msg := equals(s.Dump(), matchingBytes[1:]); msg != "" {
			t.Fatal(msg)
		}
		s.Pop()
		if msg := equals(s.Dump(), matchingBytes); msg != "" {
			t.Fatal(msg)
		}

		s.Push()
		if msg := noerror(s.Want(6)); msg != "" {
			t.Fatal(msg)
		}
		s.Advance()
		if msg := equals(s.Dump(), matchingBytes[6:]); msg != "" {
			t.Fatal(msg)
		}
		s.Drop()
		if msg := equals(s.Dump(), matchingBytes[6:]); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkStateStack(b *testing.B) {
	s := pars.FromString("Hello world!")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Push()
		s.Pop()
	}
}
