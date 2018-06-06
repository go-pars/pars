package examples

import (
	"strings"
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/require"
)

func TestPolish(t *testing.T) {
	t.Run("matches number", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("42"))
		result, err := pars.Apply(Expression, s)
		require.NoError(t, err)
		require.Equal(t, 42.0, result)
	})

	t.Run("matches flat operation", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("+ 2 2"))
		result, err := pars.Apply(Expression, s)
		require.NoError(t, err)
		require.Equal(t, 4.0, result)
	})

	t.Run("matches nested operation", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("* - 5 6 7"))
		result, err := pars.Apply(Expression, s)
		require.NoError(t, err)
		require.Equal(t, -7.0, result)
	})

	t.Run("matches nested operation", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("- 5 * 6 7"))
		result, err := pars.Apply(Expression, s)
		require.NoError(t, err)
		require.Equal(t, -37.0, result)
	})
}
