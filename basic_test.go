package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/assert"
)

func TestEpsilon(t *testing.T) {
	s := pars.FromString("Hello world!")
	err := pars.Epsilon(s, pars.Void)
	assert.NoError(t, err)
}

func BenchmarkEpsilon(b *testing.B) {
	s := pars.FromString("Hello world!")
	p := pars.Epsilon
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.Void)
	}
}

func TestFail(t *testing.T) {
	s := pars.FromString("Hello world!")
	err := pars.Fail(s, pars.Void)
	assert.Error(t, err)
}

func BenchmarkFail(b *testing.B) {
	s := pars.FromString("Hello world!")
	p := pars.Fail
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.Void)
	}
}

func TestHead(t *testing.T) {
	s := pars.FromString("Hello world!")

	t.Run("matches at head", func(t *testing.T) {
		err := pars.Head(s, pars.Void)
		assert.NoError(t, err)
	})

	assert.NoError(t, s.Want(1))
	s.Advance()

	t.Run("fails otherwise", func(t *testing.T) {
		err := pars.Head(s, pars.Void)
		assert.Error(t, err)
	})
}

func BenchmarkHead(b *testing.B) {
	s := pars.FromString("Hello world!")
	p := pars.Head

	b.Run("is head", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p(s, pars.Void)
		}
	})

	assert.NoError(b, s.Want(1))
	s.Advance()

	b.Run("not head", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			p(s, pars.Void)
		}
	})
}

func TestEOF(t *testing.T) {
	s := pars.FromString("Hello world!")

	t.Run("fails if not at EOF", func(t *testing.T) {
		err := pars.EOF(s, pars.Void)
		assert.Error(t, err)
	})

	for s.Want(1) == nil {
		s.Advance()
	}

	t.Run("matches otherwise", func(t *testing.T) {
		err := pars.EOF(s, pars.Void)
		assert.NoError(t, err)
	})
}

func BenchmarkEOF(b *testing.B) {
	s := pars.FromString("Hello world!")
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
