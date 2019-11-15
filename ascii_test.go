package pars_test

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/ktnyt/ascii"
	"github.com/ktnyt/assert"
	"github.com/ktnyt/bench"
	"github.com/ktnyt/pars"
)

func TestFilter(t *testing.T) {
	validTests := []struct {
		Set    []byte
		Filter ascii.Filter
	}{
		{ascii.Null, ascii.IsNull},
		{ascii.Graphic, ascii.IsGraphic},
		{ascii.Control, ascii.IsControl},
		{ascii.Space, ascii.IsSpace},
		{ascii.Upper, ascii.IsUpper},
		{ascii.Lower, ascii.IsLower},
		{ascii.Letter, ascii.IsLetter},
		{ascii.Digit, ascii.IsDigit},
		{ascii.Latin, ascii.IsLatin},
	}

	validCases := make([]assert.F, len(validTests))
	for i, tc := range validTests {
		p := pars.Exact(pars.Many(pars.Filter(tc.Filter)).Map(pars.Cat))
		e := pars.NewTokenResult(tc.Set)
		validCases[i] = MatchingCase(p, tc.Set, e, len(tc.Set))
	}

	assert.Apply(t,
		assert.C("matching", validCases...),
	)
}

func BenchmarkFilter(b *testing.B) {
	validTests := []struct {
		Set    []byte
		Filter ascii.Filter
	}{
		{ascii.Null, ascii.IsNull},
		{ascii.Graphic, ascii.IsGraphic},
		{ascii.Control, ascii.IsControl},
		{ascii.Space, ascii.IsSpace},
		{ascii.Upper, ascii.IsUpper},
		{ascii.Lower, ascii.IsLower},
		{ascii.Letter, ascii.IsLetter},
		{ascii.Digit, ascii.IsDigit},
		{ascii.Latin, ascii.IsLatin},
	}

	validCases := make([]bench.F, len(validTests))
	for i, bc := range validTests {
		v := reflect.ValueOf(bc.Filter)
		f := runtime.FuncForPC(v.Pointer())
		name := f.Name()

		p := pars.Filter(bc.Filter)
		validCases[i] = bench.C(name, benchmark(p, bc.Set))
	}

	bench.Apply(b, validCases...)
}
