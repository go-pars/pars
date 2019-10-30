package pars

import (
	"bytes"
	"fmt"
	"strings"
)

var reptbl = []string{
	"nul", "soh", "stx", "etx", "eot", "enq", "ack", "bel",
	"bs", "ht", "nl", "vt", "np", "cr", "so", "si",
	"dle", "dc1", "dc2", "dc3", "dc4", "nak", "syn", "etb",
	"can", "em", "sub", "esc", "fs", "gs", "rs", "us",
	"sp", "!", "\"", "#", "$", "%", "&", "'",
	"(", ")", "*", "+", " , ", "-", ".", "/",
	"0", "1", "2", "3", "4", "5", "6", "7",
	"8", "9", ":", ";", "<", "=", ">", "?",
	"@", "A", "B", "C", "D", "E", "F", "G",
	"H", "I", "J", "K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T", "U", "V", "W",
	"X", "Y", "Z", "[", "\\", "]", "^", "_",
	"`", "a", "b", "c", "d", "e", "f", "g",
	"h", "i", "j", "k", "l", "m", "n", "o",
	"p", "q", "r", "s", "t", "u", "v", "w",
	"x", "y", "z", "{", "|", "}", "~", "del",
}

func rep(c byte) string {
	if int(c) < len(reptbl) {
		return fmt.Sprintf("`%s`", reptbl[int(c)])
	}
	return fmt.Sprintf("0x%x", c)
}

func reps(p []byte) []string {
	r := make([]string, len(p))
	for i, c := range p {
		r[i] = rep(c)
	}
	return r
}

// Byte will match the given byte.
func Byte(c byte) Parser {
	name := fmt.Sprintf("Byte(%s)", rep(c))
	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return NewTraceError(name, err)
		}
		if state.Head() != c {
			return NewMismatchError(name, c, state.Position())
		}
		result.SetToken(asBytes(c))
		state.Advance()
		return nil
	}
}

// AnyByte will match any of the given bytes.
func AnyByte(p ...byte) Parser {
	switch len(p) {
	case 0:
		return Epsilon
	case 1:
		return Byte(p[0])
	default:
		s := string(p)
		mismatch := func(c byte) bool { return strings.IndexByte(s, c) < 0 }

		name := fmt.Sprintf("AnyByte(%s)", strings.Join(reps(p), ", "))

		return func(state *State, result *Result) error {
			if err := state.Want(1); err != nil {
				return NewTraceError(name, err)
			}
			c := state.Head()
			if mismatch(c) {
				return NewMismatchError(name, s, state.Position())
			}
			result.SetToken(asBytes(c))
			state.Advance()
			return nil
		}
	}
}

func sign(i int) int {
	if i > 0 {
		return 1
	}
	if i < 0 {
		return -1
	}
	return 0
}

// ByteRange will match any byte within the given range.
func ByteRange(begin, end byte) Parser {
	switch sign(int(end - begin)) {
	case -1:
		panic(fmt.Errorf("byte `%s` is greater than `%s`", rep(begin), rep(end)))
	case 0:
		return Byte(begin)
	default:
		name := fmt.Sprintf("ByteRange(%s, %s)", rep(begin), rep(end))
		e := fmt.Sprintf("in range %s-%s", rep(begin), rep(end))

		return func(state *State, result *Result) error {
			if err := state.Want(1); err != nil {
				return NewTraceError(name, err)
			}
			c := state.Head()
			if c < begin || end < c {
				return NewMismatchError(name, e, state.Position())
			}
			result.SetToken(asBytes(c))
			state.Advance()
			return nil
		}
	}
}

// Bytes will match the given sequence of bytes.
func Bytes(p []byte) Parser {
	n := len(p)
	switch n {
	case 0:
		return Epsilon
	case 1:
		return Byte(p[0])
	default:
		e := fmt.Sprintf("[%s]", strings.Join(reps(p), ", "))
		name := fmt.Sprintf("Bytes([%s])", e)
		return func(state *State, result *Result) error {
			if err := state.Want(n); err != nil {
				return NewTraceError(name, err)
			}
			buffer := state.Buffer()[:n]
			if !bytes.Equal(p, buffer) {
				return NewMismatchError(name, e, state.Position())
			}
			result.SetToken(buffer)
			state.Advance()
			return nil
		}
	}
}
