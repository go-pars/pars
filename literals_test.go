package pars_test

import (
	"strings"
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/require"
)

func TestQuoted(t *testing.T) {
	r := pars.Result{}

	t.Run("matches quoted", func(t *testing.T) {
		e := `some "string"`
		s := pars.NewState(strings.NewReader(`"some \"string\""`))
		err := pars.Quoted('"')(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})
}

func BenchmarkQuoted(b *testing.B) {
	p := pars.Dry(pars.Quoted('"'))
	s := pars.NewState(strings.NewReader(`"some \"string\""`))
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestInteger(t *testing.T) {
	r := pars.Result{}

	t.Run("matches zero", func(t *testing.T) {
		e := "0"
		s := pars.NewState(strings.NewReader("0123"))
		err := pars.Integer(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("matches integer", func(t *testing.T) {
		e := "42"
		s := pars.NewState(strings.NewReader("42: the answer"))
		err := pars.Integer(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("matches negative", func(t *testing.T) {
		e := "-42"
		s := pars.NewState(strings.NewReader("-42: negative"))
		err := pars.Integer(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("a string"))
		err := pars.Integer(s, pars.VoidResult)
		require.Error(t, err)
	})
}

func BenchmarkInteger(b *testing.B) {
	p := pars.Dry(pars.Integer)
	s := pars.NewState(strings.NewReader("42: the answer"))
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestNumber(t *testing.T) {
	r := pars.Result{}

	t.Run("matches zero", func(t *testing.T) {
		e := "0"
		s := pars.NewState(strings.NewReader("0123"))
		err := pars.Number(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("matches integer", func(t *testing.T) {
		e := "42"
		s := pars.NewState(strings.NewReader("42: the answer"))
		err := pars.Number(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("matches negative", func(t *testing.T) {
		e := "-42"
		s := pars.NewState(strings.NewReader("-42: negative"))
		err := pars.Number(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("matches decimal", func(t *testing.T) {
		e := "42.42"
		s := pars.NewState(strings.NewReader("42.42: a decimal"))
		err := pars.Number(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("matches exponent", func(t *testing.T) {
		e := "42e42"
		s := pars.NewState(strings.NewReader("42e42: an exponent"))
		err := pars.Number(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("matches decimal and exponent", func(t *testing.T) {
		e := "42.42e42"
		s := pars.NewState(strings.NewReader("42.42e42: a number"))
		err := pars.Number(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})
}

func BenchmarkNumber(b *testing.B) {
	p := pars.Dry(pars.Number)
	s := pars.NewState(strings.NewReader("42.42e42: the answer"))
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}
