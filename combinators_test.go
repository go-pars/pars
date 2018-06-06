package pars_test

import (
	"strings"
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/require"
)

func TestSeq(t *testing.T) {
	e := []rune{'🍺', '🍣'}
	q := make([]pars.ParserLike, len(e))
	for i := range e {
		q[i] = e[i]
	}
	p := pars.Seq(q...)

	t.Run("matches sequence", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("🍺🍣"))
		r := pars.Result{}
		err := p(s, &r)
		require.NoError(t, err)
		var c []pars.Result
		require.IsType(t, c, r.Children)
		c = r.Children
		require.Equal(t, len(e), len(c))
		for i := range e {
			require.IsType(t, e[i], c[i].Value)
			require.Equal(t, e[i], c[i].Value)
		}
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("🍺🍖"))
		r := pars.Result{}
		err := p(s, &r)
		require.Error(t, err)
	})
}

func BenchmarkByteSeq(b *testing.B) {
	p := pars.Dry(pars.Seq(pars.Byte('4'), pars.Byte('2')))
	s := pars.NewState(strings.NewReader("42"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func BenchmarkRuneSeq(b *testing.B) {
	p := pars.Dry(pars.Seq('🍺', '🍣'))
	s := pars.NewState(strings.NewReader("🍺🍣"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestAny(t *testing.T) {
	r := pars.Result{}

	t.Run("matches any", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("hello world"))
		err := pars.Any("hello", "world")(s, &r)
		require.NoError(t, err)
		require.IsType(t, "hello", r.Value)
		require.Equal(t, "hello", r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("hello world"))
		err := pars.Any("nihao", "world")(s, &r)
		require.Error(t, err)
	})
}

func BenchmarkAny(b *testing.B) {
	p := pars.Dry(pars.Any('🍺', '🍣'))
	s := pars.NewState(strings.NewReader("🍣🍺"))
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestTry(t *testing.T) {
	r := pars.Result{}

	t.Run("try to match", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("🍺🍣"))
		err := pars.Try('🍺')(s, &r)
		require.NoError(t, err)
		require.IsType(t, '🍺', r.Value)
		require.Equal(t, '🍺', r.Value)
	})

	t.Run("returns no error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("🍖🍣"))
		err := pars.Try('🍺')(s, &r)
		require.NoError(t, err)
	})
}

func BenchmarkTry(b *testing.B) {
	p := pars.Dry(pars.Try('🍺'))
	s := pars.NewState(strings.NewReader("🍺🍣"))
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestMany(t *testing.T) {
	r := pars.Result{}

	t.Run("matches many", func(t *testing.T) {
		e := '🍺'
		s := pars.NewState(strings.NewReader("🍺🍺🍺"))
		err := pars.Many('🍺', 0)(s, &r)
		require.NoError(t, err)
		var c []pars.Result
		require.IsType(t, c, r.Children)
		for i := range c {
			require.IsType(t, e, c[i].Value)
			require.Equal(t, e, c[i].Value)
		}
	})
}

func BenchmarkMany(b *testing.B) {
	p := pars.Dry(pars.Many('🍺', 0))
	s := pars.NewState(strings.NewReader("🍺🍺🍺"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}
