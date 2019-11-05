package pars_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ktnyt/pars"
)

func assert(condition bool) string {
	if !condition {
		return "expected true"
	}
	return ""
}

func noerror(err error) string {
	if err != nil {
		return fmt.Sprintf("unexpected error: %s", err.Error())
	}
	return ""
}

func iserror(err error) string {
	if err == nil {
		return "expected error"
	}
	return ""
}

func equals(act, exp interface{}) string {
	if !reflect.DeepEqual(exp, act) {
		return fmt.Sprintf("expected: %#v\n  actual: %#v", exp, act)
	}
	return ""
}

func try(msgs ...string) string {
	for _, msg := range msgs {
		if msg != "" {
			return msg
		}
	}
	return ""
}

var matchingString = "Hello world!"
var matchingBytes = []byte(matchingString)
var matchingRunes = []rune(matchingString)

var mismatchString = "Goodbye world!"
var mismatchBytes = []byte(mismatchString)
var mismatchRunes = []rune(mismatchString)

func matching(p pars.Parser, in []byte, e *pars.Result, out []byte) string {
	s := pars.FromBytes(in)
	r := &pars.Result{}
	return try(
		noerror(p(s, r)),
		equals(r.Token, e.Token),
		equals(r.Value, e.Value),
		equals(r.Children, e.Children),
		equals(s.Dump(), out),
	)
}

func mismatch(p pars.Parser, in []byte) string {
	s := pars.FromBytes(in)
	r := &pars.Result{}
	return try(
		iserror(p(s, r)),
		equals(r, &pars.Result{}),
		equals(s.Dump(), in),
	)
}

type benchFunc func(*testing.B)

type benchCase struct {
	Name string
	Func benchFunc
}

func benchmark(p pars.Parser, in []byte) benchFunc {
	s := pars.FromBytes(in)
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Push()
			p(s, pars.Void)
			s.Pop()
		}
	}
}

func combineBench(cases ...benchCase) benchFunc {
	return func(b *testing.B) {
		for _, bc := range cases {
			b.Run(bc.Name, bc.Func)
		}
	}
}
