package pars_test

import (
	"strings"
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/require"
)

func TestDelim(t *testing.T) {
	r := pars.Result{}

	t.Run("matches delimited sequence", func(t *testing.T) {
		e := []rune{'🍺', '🍺', '🍺'}
		s := pars.NewState(strings.NewReader("🍺🍣🍺🍣🍺🍣"))
		err := pars.Delim('🍺', '🍣')(s, &r)
		require.NoError(t, err)
		c := r.Children
		require.Equal(t, len(e), len(c))
		for i := range e {
			require.IsType(t, e[i], c[i].Value)
			require.Equal(t, e[i], c[i].Value)
		}
	})
}

func BenchmarkDelim(b *testing.B) {
	p := pars.Dry(pars.Delim('🍺', '🍣'))
	s := pars.NewState(strings.NewReader("🍺🍣🍺🍣🍺🍣"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestSep(t *testing.T) {
	r := pars.Result{}

	t.Run("matches separated sequence", func(t *testing.T) {
		e := []rune{'🍺', '🍺', '🍺'}
		s := pars.NewState(strings.NewReader("🍺🍣 🍺🍣 🍺🍣"))
		err := pars.Sep('🍺', '🍣')(s, &r)
		require.NoError(t, err)
		c := r.Children
		require.Equal(t, len(e), len(c))
		for i := range e {
			require.IsType(t, e[i], c[i].Value)
			require.Equal(t, e[i], c[i].Value)
		}
	})
}

func BenchmarkSep(b *testing.B) {
	p := pars.Dry(pars.Sep('🍺', '🍣'))
	s := pars.NewState(strings.NewReader("🍺🍣 🍺🍣 🍺🍣"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestPhrase(t *testing.T) {
	e := []rune{'🍺', '🍣'}
	q := make([]pars.ParserLike, len(e))
	for i := range e {
		q[i] = e[i]
	}
	p := pars.Phrase(q...)

	t.Run("matches phrase", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("🍺 🍣"))
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
		s := pars.NewState(strings.NewReader("🍺 🍖"))
		r := pars.Result{}
		err := p(s, &r)
		require.Error(t, err)
	})
}

func BenchmarkPhrase(b *testing.B) {
	p := pars.Dry(pars.Phrase(pars.Byte('4'), pars.Byte('2')))
	s := pars.NewState(strings.NewReader("4 2"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestUntil(t *testing.T) {
	p := pars.Until('🍺')

	t.Run("matches until", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("🍣🍖🍺"))
		r := pars.Result{}
		err := p(s, &r)
		s.Clear()
		require.NoError(t, err)
		require.Equal(t, []byte("🍣🍖"), r.Value)
		require.Equal(t, []byte("🍺"), s.Buffer)
	})

	t.Run("returns no error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("🍺🍺🍺"))
		err := p(s, pars.VoidResult)
		require.NoError(t, err)
	})
}

func BenchmarkUntil(b *testing.B) {
	p := pars.Dry(pars.Until('🍺'))
	s := pars.NewState(strings.NewReader("🍣🍖🍺"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}
