package pars_test

import (
	"strings"
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/assert"
)

func TestByte(t *testing.T) {
	t.Run("matches byte", func(t *testing.T) {
		e := byte('H')
		s := pars.FromString("Hello world!")
		r := pars.Result{}
		err := pars.Byte(e)(s, &r)
		assert.NoError(t, err)
		assert.IsType(t, e, r.Token[0])
		assert.Equal(t, e, r.Token[0])
		assert.Nil(t, r.Value)
		assert.Nil(t, r.Children)
	})

	t.Run("returns error", func(t *testing.T) {
		e := byte('h')
		s := pars.FromString("Hello world!")
		r := pars.Result{}
		err := pars.Byte(e)(s, &r)
		assert.Error(t, err)
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
