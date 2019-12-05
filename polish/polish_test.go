package examples

import (
	"testing"

	"gopkg.in/ktnyt/assert.v1"
	"gopkg.in/ktnyt/pars.v2"
)

func testCase(s string, e float64) assert.F {
	r, err := Expression.Parse(pars.FromString(s))
	return assert.All(assert.NoError(err), assert.Equal(r.Value, e))
}

func TestPolish(t *testing.T) {
	assert.Apply(t,
		assert.C("matches number", testCase("42", 42.0)),
		assert.C("matches flat operation", testCase("+ 2 2", 4.0)),
		assert.C("matches nested operation", testCase("* - 5 6 7", -7.0)),
		assert.C("matches nested operation", testCase("- 5 * 6 7", -37.0)),
	)
}
