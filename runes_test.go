package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRune(t *testing.T) {
	t.Run("no argument", func(t *testing.T) {
		e := rune('H')
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.NoError(t, pars.Rune()(s, &r))
		require.NotEmpty(t, r.Value)
		assert.Nil(t, r.Token)
		assert.Equal(t, r.Value, e)
		assert.Nil(t, r.Children)
	})

	t.Run("single argument", func(t *testing.T) {
		t.Run("matching", func(t *testing.T) {
			e := rune('H')
			s := pars.FromString("Hello world!")
			r := pars.Result{}

			require.NoError(t, pars.Rune(e)(s, &r))
			require.NotEmpty(t, r.Value)
			assert.Nil(t, r.Token)
			assert.Equal(t, r.Value, e)
			assert.Nil(t, r.Children)
		})

		t.Run("mismatch", func(t *testing.T) {
			s := pars.FromString("Hello world!")
			r := pars.Result{}

			require.Error(t, pars.Rune('h')(s, &r))
			assert.Nil(t, r.Token)
			assert.Nil(t, r.Value)
			assert.Nil(t, r.Children)
		})
	})

	t.Run("many arguments", func(t *testing.T) {
		t.Run("matching", func(t *testing.T) {
			e := rune('H')
			s := pars.FromString("Hello world!")
			r := pars.Result{}

			require.NoError(t, pars.Rune('h', e)(s, &r))
			require.NotEmpty(t, r.Value)
			assert.Nil(t, r.Token)
			assert.Equal(t, r.Value, e)
			assert.Nil(t, r.Children)
		})

		t.Run("mismatch", func(t *testing.T) {
			s := pars.FromString("Hello world!")
			r := pars.Result{}

			assert.Error(t, pars.Rune('h', 'w')(s, &r))
			assert.Nil(t, r.Token)
			assert.Nil(t, r.Value)
			assert.Nil(t, r.Children)
		})
	})
}

func BenchmarkRune(b *testing.B) {
	s := pars.FromString("Hello world!")

	b.Run("no argument", func(b *testing.B) {
		p := pars.Rune()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("single argument", func(b *testing.B) {
		b.Run("matching", func(b *testing.B) {
			p := pars.Rune('H')
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Push()
				p(s, pars.Void)
				s.Pop()
			}
		})

		b.Run("mismatch", func(b *testing.B) {
			p := pars.Rune('h')
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Push()
				p(s, pars.Void)
				s.Pop()
			}
		})
	})

	b.Run("many arguments", func(b *testing.B) {
		b.Run("matching first", func(b *testing.B) {
			p := pars.Rune('H', 'h')
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Push()
				p(s, pars.Void)
				s.Pop()
			}
		})

		b.Run("matching second", func(b *testing.B) {
			p := pars.Rune('h', 'H')
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Push()
				p(s, pars.Void)
				s.Pop()
			}
		})

		b.Run("mismatch", func(b *testing.B) {
			p := pars.Rune('h', 'w')
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Push()
				p(s, pars.Void)
				s.Pop()
			}
		})
	})
}
