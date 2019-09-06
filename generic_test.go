package pars_test

import (
	"strings"
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/require"
)

func TestLine(t *testing.T) {
	r := pars.Result{}

	t.Run("matches line", func(t *testing.T) {
		e := "this is a line"
		b := "and this is another"
		s := pars.NewState(strings.NewReader("this is a line\nand this is another"))
		err := pars.Line(s, &r)
		pars.Cut(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
		require.Equal(t, b, string(s.Buffer))
	})
}

func BenchmarkLine(b *testing.B) {
	p := pars.Dry(pars.Line)
	s := pars.NewState(strings.NewReader("this is a line\nand this is another"))
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestWord(t *testing.T) {
	r := pars.Result{}

	t.Run("matches word", func(t *testing.T) {
		e := "hello"
		b := " world"
		s := pars.NewState(strings.NewReader("hello world"))
		err := pars.Word(s, &r)
		pars.Cut(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
		require.Equal(t, b, string(s.Buffer))
	})
}

func BenchmarkWord(b *testing.B) {
	p := pars.Dry(pars.Word)
	s := pars.NewState(strings.NewReader("hello world"))
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}
