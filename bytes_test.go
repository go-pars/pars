package pars_test

import (
	"strings"
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestByte(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		e := byte('H')
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.NoError(t, pars.Byte(e)(s, &r))
		require.NotEmpty(t, r.Token)
		assert.Equal(t, r.Token[0], e)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})

	t.Run("mismatch", func(t *testing.T) {
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.Error(t, pars.Byte('h')(s, &r))
		assert.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})
}

func BenchmarkByte(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))

	b.Run("matching", func(b *testing.B) {
		p := pars.Byte('H')
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("mismatch", func(b *testing.B) {
		p := pars.Byte('h')
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})
}

func TestAnyByte(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		e := byte('H')
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.NoError(t, pars.AnyByte('h', e)(s, &r))
		require.NotEmpty(t, r.Token)
		assert.Equal(t, r.Token[0], e)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})

	t.Run("mismatch", func(t *testing.T) {
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		assert.Error(t, pars.AnyByte('h', 'w')(s, &r))
		assert.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})
}

func BenchmarkAnyByte(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))

	b.Run("matching first", func(b *testing.B) {
		p := pars.AnyByte('H', 'h')
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("matching second", func(b *testing.B) {
		p := pars.AnyByte('h', 'H')
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("mismatch", func(b *testing.B) {
		p := pars.AnyByte('h', 'w')
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})
}

func TestByteRange(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		e := byte('H')
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.NoError(t, pars.ByteRange('A', 'Z')(s, &r))
		require.NotEmpty(t, r.Token)
		assert.Equal(t, r.Token[0], e)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})

	t.Run("mismatch", func(t *testing.T) {
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.Error(t, pars.ByteRange('a', 'z')(s, &r))
		require.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})
}

func BenchmarkRangeByte(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))

	b.Run("matching", func(b *testing.B) {
		p := pars.ByteRange('A', 'Z')
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})

	b.Run("mismatch", func(b *testing.B) {
		p := pars.ByteRange('a', 'z')
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	})
}

func TestBytes(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		e := []byte("Hello")
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.NoError(t, pars.Bytes(e)(s, &r))
		require.NotEmpty(t, r.Token)
		assert.ElementsMatch(t, r.Token, e)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})

	t.Run("mismatch", func(t *testing.T) {
		s := pars.FromString("Hello world!")
		r := pars.Result{}

		require.Error(t, pars.Bytes([]byte("hello"))(s, &r))
		require.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})
}
