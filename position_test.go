package pars

import "testing"

var positionLessTests = []struct {
	a, b Position
	ok   bool
}{
	{Position{0, 0}, Position{0, 0}, false},
	{Position{0, 0}, Position{0, 1}, true},
	{Position{0, 0}, Position{1, 0}, true},
	{Position{1, 0}, Position{0, 0}, false},
}

func TestPositionLess(t *testing.T) {
	for _, tt := range positionLessTests {
		ok := tt.a.Less(tt.b)
		if ok != tt.ok {
			t.Errorf("<%s>.Less(<%s>) = %t, want %t", tt.a, tt.b, ok, tt.ok)
		}
	}
}
