package examples

import (
	"testing"

	"github.com/go-pars/pars"
)

var polishTestCases = []struct {
	in  string
	out float64
}{
	{"42", 42.0},
	{"+ 2 2", 4.0},
	{"* - 5 6 7", -7.0},
	{"- 5 * 6 7", -37.0},
}

func TestPolish(t *testing.T) {
	for _, tt := range polishTestCases {
		out, err := Expression.Parse(pars.FromString(tt.in))
		if err != nil {
			t.Errorf("Expresson.Parse(%q): %v", tt.in, err)
			return
		}
		f := out.Value.(float64)
		if f != tt.out {
			t.Errorf("Expression.Parse(%q) = %f, wanted %f", tt.in, f, tt.out)
		}
	}
}
