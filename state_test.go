package pars_test

import (
	"io"
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	t.Run("Read", func(t *testing.T) {
		e := "Hello world!"
		s := pars.FromString(e)
		p := make([]byte, len(e))
		n, err := s.Read(p)

		assert.Equal(t, n, len(e))
		assert.NoError(t, err)
		assert.Equal(t, string(p), e)
		assert.Empty(t, s.Dump())
	})

	t.Run("Want", func(t *testing.T) {
		e := "Hello world!"
		s := pars.FromString(e)

		assert.NoError(t, s.Want(len(e)))
		assert.Equal(t, s.Want(len(e)+1), io.EOF)
	})

	t.Run("Advance", func(t *testing.T) {
		e := "Hello world!"
		s := pars.FromString(e)

		assert.NoError(t, s.Want(1))
		s.Advance()
		assert.Equal(t, string(s.Dump()), e[1:])

		assert.NoError(t, s.Want(5))
		s.Advance()
		assert.Equal(t, string(s.Dump()), e[6:])
	})

	t.Run("Stack", func(t *testing.T) {
		e := "Hello world!"
		s := pars.FromString(e)

		s.Push()
		assert.NoError(t, s.Want(1))
		s.Advance()
		assert.Equal(t, string(s.Dump()), e[1:])
		s.Push()
		assert.NoError(t, s.Want(5))
		s.Advance()
		assert.Equal(t, string(s.Dump()), e[6:])
		s.Pop()
		assert.Equal(t, string(s.Dump()), e[1:])
		s.Pop()
		assert.Equal(t, string(s.Dump()), e)

		s.Push()
		assert.NoError(t, s.Want(6))
		s.Advance()
		assert.Equal(t, string(s.Dump()), e[6:])
		s.Drop()
		assert.Equal(t, string(s.Dump()), e[6:])
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
