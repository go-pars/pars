package pars

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

func runeRep(r rune) string {
	n := utf8.RuneLen(r)
	if n > 1 {
		return fmt.Sprintf("%c", r)
	}
	p := make([]byte, 1)
	utf8.EncodeRune(p, r)
	return byteRep(p[0])
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
		if err := state.Want(i + 1); err != nil {
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

func Rune(rs ...rune) Parser {
	switch len(rs) {
	case 0:
		return func(state *State, result *Result) error {
			r, err := readRune(state)
			if err != nil {
				return NewTraceError("Rune()", err)
			}
			result.SetValue(r)
			state.Advance()
			return nil
		}
	case 1:
		e := rs[0]
		rep := runeRep(e)
		name := fmt.Sprintf("Rune(%s)", rep)

		return func(state *State, result *Result) error {
			r, err := readRune(state)
			if err != nil {
				return NewTraceError(name, err)
			}
			if r != e {
				return NewMismatchError(name, rep, state.Position())
			}
			result.SetValue(r)
			state.Advance()
			return nil
		}
	default:
		reps := strings.Join(runeReps(rs), ", ")
		name := fmt.Sprintf("Rune(%s)", reps)

		s := string(rs)
		mismatch := func(r rune) bool { return !strings.ContainsRune(s, r) }

		return func(state *State, result *Result) error {
			r, err := readRune(state)
			if err != nil {
				return NewTraceError(name, err)
			}
			if mismatch(r) {
				return NewMismatchError(name, reps, state.Position())
			}
			result.SetValue(r)
			state.Advance()
			return nil
		}
	}
}
