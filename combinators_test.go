package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeq(t *testing.T) {
	e := []rune("Hello")
	q := make([]pars.ParserLike, len(e))
	for i := range e {
		q[i] = e[i]
	}
	p := pars.Seq(q...)

	t.Run("matching", func(t *testing.T) {
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.NoError(t, p(s, &r))
		assert.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		require.NotEmpty(t, r.Children)
		for i, child := range r.Children {
			assert.Equal(t, e[i], child.Value)
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		s := pars.FromString("hello world!")
		r := pars.Result{}

		require.Error(t, p(s, &r))
		assert.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})
}

func BenchmarkSeq(b *testing.B) {
	e := []rune("Hello")
	q := make([]pars.ParserLike, len(e))
	for i := range e {
		q[i] = e[i]
	}
	p := pars.Seq(q...)

	b.Run("matching", func(b *testing.B) {
		s := pars.FromString("Hello world!")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("mismatch", func(b *testing.B) {
		s := pars.FromString("hello world!")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})
}

func TestAny(t *testing.T) {
	first := "Hello"
	second := "Goodbye"
	p := pars.Any(first, second)

	t.Run("matching first", func(t *testing.T) {
		s := pars.FromString(first + " world!")
		r := pars.Result{}

		require.NoError(t, p(s, &r))
		assert.Nil(t, r.Token)
		assert.Equal(t, r.Value, first)
		assert.Nil(t, r.Children)
	})

	t.Run("matching second", func(t *testing.T) {
		s := pars.FromString(second + " world!")
		r := pars.Result{}

		require.NoError(t, p(s, &r))
		assert.Nil(t, r.Token)
		assert.Equal(t, r.Value, second)
		assert.Nil(t, r.Children)
	})

	t.Run("mismatch", func(t *testing.T) {
		s := pars.FromString("Nihao world!")
		r := pars.Result{}

		require.Error(t, p(s, &r))
		assert.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})
}

func BenchmarkAny(b *testing.B) {
	first := "Hello"
	second := "Goodbye"
	p := pars.Any(first, second)

	b.Run("matching first", func(b *testing.B) {
		s := pars.FromString(first + " world!")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("matching second", func(b *testing.B) {
		s := pars.FromString(second + " world!")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("mismatch", func(b *testing.B) {
		s := pars.FromString("Nihao world!")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})
}

func TestMaybe(t *testing.T) {
	e := "Hello"
	p := pars.Maybe(e)

	t.Run("matching", func(t *testing.T) {
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.NoError(t, p(s, &r))
		assert.Nil(t, r.Token)
		assert.Equal(t, r.Value, e)
		assert.Nil(t, r.Children)
	})

	t.Run("mismatch", func(t *testing.T) {
		s := pars.FromString("hello world!")
		r := pars.Result{}

		require.NoError(t, p(s, &r))
		assert.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})
}

func BenchmarkMaybe(b *testing.B) {
	p := pars.Maybe("Hello")

	b.Run("matching", func(b *testing.B) {
		s := pars.FromString("Hello world!")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("mismatch", func(b *testing.B) {
		s := pars.FromString("hello world!")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})
}
