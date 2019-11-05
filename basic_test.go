package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
)

func TestEpsilon(t *testing.T) {
	s := pars.FromString(matchingString)
	if msg := noerror(pars.Epsilon(s, pars.Void)); msg != "" {
		t.Fatal(msg)
	}
}

func BenchmarkEpsilon(b *testing.B) {
	s := pars.FromString(matchingString)
	p := pars.Epsilon
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.Void)
	}
}

func TestFail(t *testing.T) {
	s := pars.FromString(matchingString)
	if msg := iserror(pars.Fail(s, pars.Void)); msg != "" {
		t.Fatal(msg)
	}
}

func BenchmarkFail(b *testing.B) {
	s := pars.FromString(matchingString)
	p := pars.Fail
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.Void)
	}
}

func TestHead(t *testing.T) {
	s := pars.FromString(matchingString)

	t.Run("matches at head", func(t *testing.T) {
		if msg := noerror(pars.Head(s, pars.Void)); msg != "" {
			t.Fatal(msg)
		}
	})

	if msg := noerror(s.Want(1)); msg != "" {
		t.Fatal(msg)
	}
	s.Advance()

	t.Run("fails otherwise", func(t *testing.T) {
		if msg := iserror(pars.Head(s, pars.Void)); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkHead(b *testing.B) {
	s := pars.FromString(matchingString)
	p := pars.Head

	b.Run("is head", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p(s, pars.Void)
		}
	})

	if msg := noerror(s.Want(1)); msg != "" {
		b.Fatal(msg)
	}
	s.Advance()

	b.Run("not head", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p(s, pars.Void)
		}
	})
}

func TestEOF(t *testing.T) {
	s := pars.FromString(matchingString)

	t.Run("fails if not at EOF", func(t *testing.T) {
		if msg := iserror(pars.EOF(s, pars.Void)); msg != "" {
			t.Fatal(msg)
		}
	})

	for s.Want(1) == nil {
		s.Advance()
	}

	t.Run("matches otherwise", func(t *testing.T) {
		if msg := noerror(pars.EOF(s, pars.Void)); msg != "" {
			t.Fatal(msg)
		}
	})
}

func BenchmarkEOF(b *testing.B) {
	s := pars.FromString(matchingString)
	p := pars.EOF

	b.Run("not EOF", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p(s, pars.Void)
		}
	})

	for s.Want(1) == nil {
		s.Advance()
	}

	b.Run("is EOF", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p(s, pars.Void)
		}
	})
}
