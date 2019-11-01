package pars_test

import (
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/assert"
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
		s := pars.FromString("🍺🍣")
		r := pars.Result{}
		err := p(s, &r)
		require.NoError(t, err)
		assert.Nil(t, r.Token)
		assert.Nil(t, r.Value)
		require.NotEmpty(t, r.Children)
		for i, child := range r.Children {
			assert.Equal(t, e[i], child.Value)
		}
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.FromString("🍺🍖")
		r := pars.Result{}
		err := p(s, &r)
		require.Error(t, err)
	})
}
