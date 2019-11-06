package pars

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ktnyt/ascii"
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
		r := rs[0]
		rep := runeRep(r)
		name := fmt.Sprintf("Rune(%s)", rep)

		n := utf8.RuneLen(r)
		p := make([]byte, n)
		utf8.EncodeRune(p, r)

		return func(state *State, result *Result) error {
			if err := state.Want(n); err != nil {
				return NewTraceError(name, err)
			}
			if !bytes.Equal(state.Buffer(), p) {
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

// RuneRange will match any rune within the given range.
func RuneRange(begin, end rune) Parser {
	switch sign(int(end - begin)) {
	case -1:
		panic(fmt.Errorf("rune `%s` is greater than `%s`", runeRep(begin), runeRep(end)))
	case 0:
		return Rune(begin)
	default:
		name := fmt.Sprintf("RuneRange(%s, %s)", runeRep(begin), runeRep(end))
		rep := fmt.Sprintf("in range %s-%s", runeRep(begin), runeRep(end))

		return func(state *State, result *Result) error {
			r, err := readRune(state)
			if err != nil {
				return NewTraceError(name, err)
			}
			if r < begin || end < r {
				return NewMismatchError(name, rep, state.Position())
			}
			result.SetValue(r)
			state.Advance()
			return nil
		}
	}
}

// Runes will match the given sequence of runes.
func Runes(rs []rune) Parser {
	n := len(rs)
	switch n {
	case 0:
		return Epsilon
	case 1:
		return Rune(rs[0])
	default:
		reps := fmt.Sprintf("[%s]", strings.Join(runeReps(rs), ", "))
		name := fmt.Sprintf("Runes([%s])", reps)
		p := []byte(string(rs))

		return func(state *State, result *Result) error {
			if err := state.Want(len(p)); err != nil {
				return NewTraceError(name, err)
			}
			if !bytes.Equal(state.Buffer(), p) {
				return NewMismatchError(name, reps, state.Position())
			}
			result.SetValue(rs)
			state.Advance()
			return nil
		}
	}
}
