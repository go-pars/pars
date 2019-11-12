package pars_test

import (
	"io"
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/pars"
)

func TestState(t *testing.T) {
	p := []byte(hello)

	t.Run("Read", func(t *testing.T) {
		s := pars.FromBytes(p)
		q := make([]byte, len(p))
		n, err := s.Read(q)

		assert.Apply(t,
			assert.NoError(err),
			assert.Equal(n, len(q)),
			assert.Equal(q, p),
			assert.Equal(s.Dump(), []byte{}),
		)
	})

	t.Run("Request", func(t *testing.T) {
		s := pars.FromBytes(p)
		assert.Apply(t,
			assert.NoError(s.Request(len(p))),
			assert.Equal(s.Request(len(p)+1), io.EOF),
		)
	})

	t.Run("Advance", func(t *testing.T) {
		s := pars.FromBytes(p)
		advance := func() { s.Advance() }

		assert.Apply(t,
			assert.NoError(s.Request(1)),
			assert.Eval(advance),
			assert.Equal(s.Dump(), p[1:]),
			assert.NoError(s.Request(5)),
			assert.Eval(advance),
			assert.Equal(s.Dump(), p[6:]),
		)
	})

	t.Run("Stack", func(t *testing.T) {
		s := pars.FromBytes(p)
		push := func() { s.Push() }
		pop := func() { s.Pop() }
		drop := func() { s.Drop() }

		assert.Apply(t,
			assert.Eval(push),
			assert.NoError(pars.Skip(s,1)),
			assert.Equal(s.Dump(), p[1:]),

			assert.Eval(push),
			assert.NoError(pars.Skip(s,5)),
			assert.Equal(s.Dump(), p[6:]),

			assert.Eval(pop),
			assert.Equal(s.Dump(), p[1:]),

			assert.Eval(pop),
			assert.Equal(s.Dump(), p),

			assert.Eval(push),
			assert.NoError(pars.Skip(s,6)),
			assert.Equal(s.Dump(), p[6:]),

			assert.Eval(drop),
			assert.Equal(s.Dump(), p[6:]),
		)
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
