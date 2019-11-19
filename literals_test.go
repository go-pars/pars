package pars_test

import (
	"strconv"
	"testing"

	"github.com/ktnyt/assert"
	"github.com/ktnyt/bench"
	"github.com/ktnyt/pars"
)

func MatchingInt(i int) assert.F {
	s := strconv.Itoa(i)
	n := len(s)
	e := pars.NewValueResult(i)
	return MatchingCase(pars.Exact(pars.Int), []byte(s), e, n)
}

func MatchingNumber(s string) assert.F {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	n := len(s)
	e := pars.NewValueResult(f)
	return MatchingCase(pars.Exact(pars.Number), []byte(s), e, n)
}

func TestInt(t *testing.T) {
	assert.Apply(t,
		assert.C("matching",
			MatchingInt(0),
			MatchingInt(42),
			MatchingInt(-42),
		),
		assert.C("mismatch", MismatchCase(pars.Int, []byte(hello))),
	)
}

func BenchmarkInt(b *testing.B) {
	b.Run("matching 0", ParserBench(pars.Int, []byte("0 is the answer")))
	b.Run("matching 42", ParserBench(pars.Int, []byte("42 is the answer")))
	b.Run("matching -42", ParserBench(pars.Int, []byte("-42 is the answer")))
	b.Run("mismatch", ParserBench(pars.Int, []byte(hello)))
}

func TestNumber(t *testing.T) {
	validValues := []string{
		"0",
		"-0",
		"1",
		"-1",
		"0.1",
		"-0.1",
		"1234",
		"-1234",
		"12.34",
		"-12.34",
		"12E0",
		"12E1",
		"12e34",
		"12E-0",
		"12e+1",
		"12e-34",
		"-12E0",
		"-12E1",
		"-12e34",
		"-12E-0",
		"-12e+1",
		"-12e-34",
		"1.2E0",
		"1.2E1",
		"1.2e34",
		"1.2E-0",
		"1.2e+1",
		"1.2e-34",
		"-1.2E0",
		"-1.2E1",
		"-1.2e34",
		"-1.2E-0",
		"-1.2e+1",
		"-1.2e-34",
		"0E0",
		"0E1",
		"0e34",
		"0E-0",
		"0e+1",
		"0e-34",
		"-0E0",
		"-0E1",
		"-0e34",
		"-0E-0",
		"-0e+1",
		"-0e-34",
	}

	invalidValues := []string{
		"",
		"invalid",
		"1.0.1",
		"1..1",
		"-1-2",
		"012a42",
		"01.2",
		"012",
		"12E12.12",
		"1e2e3",
		"1e+-2",
		"1e--23",
		"1e",
		"e1",
		"1e+",
		"1ea",
		"1a",
		"1.a",
		"1.",
		"01",
		"1.e1",
	}

	validCases := make([]assert.F, len(validValues))
	for i, s := range validValues {
		validCases[i] = assert.C(s, MatchingNumber(s))
	}

	invalidCases := make([]assert.F, len(invalidValues))
	for i, s := range invalidValues {
		invalidCases[i] = assert.C(s, MismatchCase(pars.Exact(pars.Number), []byte(s)))
	}

	assert.Apply(t,
		assert.C("matching", validCases...),
		assert.C("mismatch", invalidCases...),
	)
}

func BenchmarkNumber(b *testing.B) {
	validValues := []string{
		"0",
		"1",
		"0.1",
		"1234",
		"12.34",
		"-61657.61667E+61673",
	}

	validCases := make([]bench.F, len(validValues))
	for i, s := range validValues {
		validCases[i] = bench.C(s, ParserBench(pars.Number, []byte(s)))
	}

	bench.Apply(b,
		bench.C("matching", validCases...),
	)
}

func TestBetween(t *testing.T) {
	p0 := []byte("(" + hello + ")")
	p1 := []byte("(\\)" + hello + ")")
	p2 := []byte("(" + hello)
	p3 := []byte("{hello, world}")
	e0 := pars.NewTokenResult([]byte(hello))
	e1 := pars.NewTokenResult([]byte("\\)" + hello))
	p := pars.Between('(', ')')

	assert.Apply(t,
		assert.C("matching",
			MatchingCase(p, p0, e0, len(p0)),
			MatchingCase(p, p1, e1, len(p1)),
		),
		assert.C("mismatch",
			MismatchCase(p, p2),
			MismatchCase(p, p3),
		),
	)
}

func BenchmarkBetween(b *testing.B) {
	p0 := []byte("(" + hello + ")")
	p1 := []byte("(\\)" + hello + ")")
	p2 := []byte("(" + hello)
	p3 := []byte("{hello, world}")
	p := pars.Between('(', ')')

	bench.Apply(b,
		bench.C("matching",
			ParserBench(p, p0),
			ParserBench(p, p1),
		),
		bench.C("mismatch",
			ParserBench(p, p2),
			ParserBench(p, p3),
		),
	)
}

func TestQuoted(t *testing.T) {
	p0 := []byte("\"" + hello + "\"")
	p1 := []byte("\"\\\"" + hello + "\"")
	p2 := []byte("'" + hello + "'")
	p3 := []byte("\"" + hello)
	e0 := pars.NewTokenResult([]byte(hello))
	e1 := pars.NewTokenResult([]byte("\\\"" + hello))
	p := pars.Quoted('"')

	assert.Apply(t,
		assert.C("matching",
			MatchingCase(p, p0, e0, len(p0)),
			MatchingCase(p, p1, e1, len(p1)),
		),
		assert.C("mismatch",
			MismatchCase(p, p2),
			MismatchCase(p, p3),
		),
	)
}

func BenchmarkQuoted(b *testing.B) {
	p0 := []byte("\"" + hello + "\"")
	p1 := []byte("\"\\\"" + hello + "\"")
	p2 := []byte("'" + hello + "'")
	p3 := []byte("\"" + hello)
	p := pars.Quoted('"')

	bench.Apply(b,
		bench.C("matching",
			ParserBench(p, p0),
			ParserBench(p, p1),
		),
		bench.C("mismatch",
			ParserBench(p, p2),
			ParserBench(p, p3),
		),
	)
}
