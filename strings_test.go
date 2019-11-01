package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		e := "Hello"
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.NoError(t, pars.String(e)(s, &r))
		require.NotEmpty(t, r.Value)
		assert.Nil(t, r.Token)
		assert.Equal(t, r.Value, e)
		assert.Nil(t, r.Children)
	})

	t.Run("mismatch", func(t *testing.T) {
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.Error(t, pars.String("hello")(s, &r))
		require.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})
}

func BenchmarkString(b *testing.B) {
	s := pars.FromString("Hello world!")

	b.Run("matching", func(b *testing.B) {
		p := pars.String("Hello")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("mismatch", func(b *testing.B) {
		p := pars.String("hello")
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})
}
