package pars_test

import (
	"reflect"
	"runtime"
	"testing"

	"gopkg.in/ktnyt/ascii.v1"
	"gopkg.in/ktnyt/assert.v1"
	"gopkg.in/ktnyt/bench.v1"
	"gopkg.in/ktnyt/pars.v2"
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
		validCases[i] = bench.C(name, ParserBench(p, bc.Set))
	}

	bench.Apply(b, validCases...)
}

func TestWord(t *testing.T) {
	p0, p1, p2 := []byte(hello), []byte(small), []byte(large)
	n := 5
	e0 := pars.NewTokenResult(p0[:n])
	e1 := pars.NewTokenResult(p1[:n])
	e2 := pars.NewTokenResult(p2[:n])
	matching := pars.Word(ascii.IsLetter)
	mismatch := pars.Word(ascii.IsSpace)

	assert.Apply(t,
		assert.C("matching",
			MatchingCase(matching, p0, e0, n),
			MatchingCase(matching, p1, e1, n),
			MatchingCase(matching, p2, e2, n),
		),
		assert.C("mismatch",
			MismatchCase(mismatch, p0),
			MismatchCase(mismatch, p1),
			MismatchCase(mismatch, p2),
		),
	)
}

func BenchmarkWord(b *testing.B) {
	p := []byte(hello)
	matching := pars.Word(ascii.IsLetter)
	mismatch := pars.Word(ascii.IsSpace)

	bench.Apply(b,
		bench.C("matching", ParserBench(matching, p)),
		bench.C("mismatch", ParserBench(mismatch, p)),
	)
}
