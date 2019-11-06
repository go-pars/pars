package pars_test

import (
	"io"
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/pars"
)

func TestEpsilon(t *testing.T) {
	p := []byte(hello)
	s := pars.FromBytes(p)
	assert.Apply(t, assert.NoError(pars.Epsilon(s, pars.Void)))
}

func BenchmarkEpsilon(b *testing.B) {
	p := []byte(hello)
	s := pars.FromBytes(p)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pars.Epsilon(s, pars.Void)
	}
}

func TestFail(t *testing.T) {
	p := []byte(hello)
	s := pars.FromBytes(p)
	assert.Apply(t, assert.IsError(pars.Fail(s, pars.Void)))
}

func BenchmarkFail(b *testing.B) {
	p := []byte(hello)
	s := pars.FromBytes(p)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pars.Fail(s, pars.Void)
	}
}

func TestHead(t *testing.T) {
	p := []byte(hello)
	s := pars.FromBytes(p)
	assert.Apply(t,
		assert.C("matches at head", assert.NoError(pars.Head(s, pars.Void))),
		assert.NoError(s.Skip(1)),
		assert.C("fails otherwise", assert.IsError(pars.Head(s, pars.Void))),
	)
}

func BenchmarkHead(b *testing.B) {
	p := []byte(hello)
	s := pars.FromBytes(p)

	b.Run("is head", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pars.Head(s, pars.Void)
		}
	})

	assert.Apply(b, assert.NoError(s.Skip(1)))

	b.Run("not head", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pars.Head(s, pars.Void)
		}
	})
}

func TestEOF(t *testing.T) {
	p := []byte(hello)
	s := pars.FromBytes(p)
	assert.Apply(t,
		assert.C("fails if not at EOF", assert.IsError(pars.EOF(s, pars.Void))),
		assert.Equal(s.Skip(len(p)+1), io.EOF),
		assert.C("matches otherwise", assert.NoError(pars.EOF(s, pars.Void))),
	)
}

func BenchmarkEOF(b *testing.B) {
	s := pars.FromString(matchingString)

	b.Run("not EOF", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pars.EOF(s, pars.Void)
		}
	})

	for s.Want(1) == nil {
		s.Advance()
	}

	b.Run("is EOF", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pars.EOF(s, pars.Void)
		}
	})
}
