package pars_test

import (
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

func TestHead(t *testing.T) {
	p := []byte(hello)
	s := pars.FromBytes(p)
	assert.Apply(t,
		assert.C("matches at head", assert.NoError(pars.Head(s, pars.Void))),
		assert.NoError(pars.Skip(s, 1)),
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

	assert.Apply(b, assert.NoError(pars.Skip(s, 1)))

	b.Run("not head", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pars.Head(s, pars.Void)
		}
	})
}

func TestEnd(t *testing.T) {
	p := []byte(hello)
	s := pars.FromBytes(p)

	assert.Apply(t,
		assert.C("fails if not at End", assert.IsError(pars.End(s, pars.Void))),
		assert.NoError(pars.Skip(s, len(p))),
		assert.C("matches otherwise", assert.NoError(pars.End(s, pars.Void))),
	)
}

func BenchmarkEnd(b *testing.B) {
	s := pars.FromString(matchingString)

	b.Run("not End", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pars.End(s, pars.Void)
		}
	})

	for s.Request(1) == nil {
		s.Advance()
	}

	b.Run("is End", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pars.End(s, pars.Void)
		}
	})
}
