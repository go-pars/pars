package pars

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	ascii "gopkg.in/ktnyt/ascii.v1"
)

func runeRep(r rune) string {
	n := utf8.RuneLen(r)
	if n > 1 {
		return fmt.Sprintf("%c", r)
	}
	p := make([]byte, 1)
	utf8.EncodeRune(p, r)
	return ascii.Rep(p[0])
}

func runeReps(p []rune) []string {
	r := make([]string, len(p))
	for i, c := range p {
		r[i] = runeRep(c)
	}
	return r
}

func readRune(state *State) (rune, error) {
	for i := 0; i < 4; i++ {
		if err := state.Request(i + 1); err != nil {
			return utf8.RuneError, err
		}
		p := state.Buffer()
		if utf8.Valid(p) {
			r, _ := utf8.DecodeRune(p)
			return r, nil
		}
	}
	return utf8.RuneError, errors.New("unable to read valid rune")
}

// Rune creates a Parser which will attempt to match the next single rune.
// If no runes are given, it will match any rune.
// Otherwise, the given runes will be tested for a match.
func Rune(rs ...rune) Parser {
	switch len(rs) {
	case 0:
		return func(state *State, result *Result) error {
			r, err := readRune(state)
			if err != nil {
				return NewNestedError("Rune", err)
			}
			result.SetValue(r)
			state.Advance()
			return nil
		}
	case 1:
		r := rs[0]
		rep := runeRep(r)
		name := fmt.Sprintf("Rune(%s)", rep)
		what := fmt.Sprintf("expected `%s`", rep)

		n := utf8.RuneLen(r)
		p := make([]byte, n)
		utf8.EncodeRune(p, r)

		return func(state *State, result *Result) error {
			if err := state.Request(n); err != nil {
				return NewNestedError(name, err)
			}
			if !bytes.Equal(state.Buffer(), p) {
				return NewError(what, state.Position())
			}
			result.SetValue(r)
			state.Advance()
			return nil
		}
	default:
		reps := strings.Join(runeReps(rs), ", ")
		name := fmt.Sprintf("Rune(%s)", reps)
		what := fmt.Sprintf("expected one of [%s]", reps)

		s := string(rs)
		mismatch := func(r rune) bool { return !strings.ContainsRune(s, r) }

		return func(state *State, result *Result) error {
			r, err := readRune(state)
			if err != nil {
				return NewNestedError(name, err)
			}
			if mismatch(r) {
				return NewError(what, state.Position())
			}
			result.SetValue(r)
			state.Advance()
			return nil
		}
	}
}

// RuneRange creates a Parser which will attempt to match any rune within the
// given range inclusively.
func RuneRange(begin, end rune) Parser {
	if begin < end {
		rbegin, rend := runeRep(begin), runeRep(end)
		name := fmt.Sprintf("RuneRange(%s, %s)", rbegin, rend)
		what := fmt.Sprintf("expected in range %s-%s", rbegin, rend)

		return func(state *State, result *Result) error {
			r, err := readRune(state)
			if err != nil {
				return NewNestedError(name, err)
			}
			if r < begin || end < r {
				return NewError(what, state.Position())
			}
			result.SetValue(r)
			state.Advance()
			return nil
		}
	}
	panic("invalid rune range")
}

// Runes creates a parser which will attempt to match the given sequence of runes.
func Runes(rs []rune) Parser {
	reps := fmt.Sprintf("[%s]", strings.Join(runeReps(rs), ", "))
	name := fmt.Sprintf("Rune(%s)", reps)
	what := fmt.Sprintf("expected [%s]", reps)
	p := []byte(string(rs))

	return func(state *State, result *Result) error {
		if err := state.Request(len(p)); err != nil {
			return NewNestedError(name, err)
		}
		if !bytes.Equal(state.Buffer(), p) {
			return NewError(what, state.Position())
		}
		result.SetValue(rs)
		state.Advance()
		return nil
	}
}
