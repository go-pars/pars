package pars

import (
	"fmt"
	"io"
	"unicode/utf8"
)

func matchByteSlice(a, b []byte) bool {
	if len(a) > len(b) {
		a, b = b, a
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Epsilon does nothing.
func Epsilon(state *State, result *Result) error {
	return nil
}

// Fail will always fail.
func Fail(state *State, result *Result) error {
	return NewParserError("must fail", state.Position)
}

// Cut the state at current position.
func Cut(state *State, result *Result) error {
	state.Clear()
	return nil
}

// Head matches if the state is at the beginning.
func Head(state *State, result *Result) error {
	if state.Position != 0 {
		return NewMismatchError("Head", []byte("Head"), state.Position)
	}
	return nil
}

// Break matches if the state is looking at the head of the buffer.
func Break(state *State, result *Result) error {
	if state.Index != 0 {
		return NewMismatchError("Break", []byte("Break"), state.Position)
	}
	return nil
}

// EOF matches if the io.Reader is at the end.
func EOF(state *State, result *Result) error {
	if err := state.Want(1); err != io.EOF {
		return NewMismatchError("EOF", []byte("EOF"), state.Position)
	}
	return nil
}

// Byte matches a given byte.
func Byte(b byte) Parser {
	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return NewTraceError("Byte", err)
		}
		if b != state.Buffer[state.Index] {
			return NewMismatchError("Byte", []byte{b}, state.Position)
		}
		result.Value = b
		state.Advance(1)
		return nil
	}
}

// NotByte matches any byte other than the given byte.
func NotByte(b byte) Parser {
	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return NewTraceError("NotByte", err)
		}
		if b == state.Buffer[state.Index] {
			return NewNotMismatchError("NotByte", []byte{b}, state.Position)
		}
		result.Value = state.Buffer[state.Index]
		state.Advance(1)
		return nil
	}
}

// Bytes matches one of the given bytes.
func Bytes(b ...byte) Parser {
	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return NewTraceError("Bytes", err)
		}
		for i := range b {
			if state.Buffer[state.Index] == b[i] {
				result.Value = state.Buffer[state.Index]
				state.Advance(1)
				return nil
			}
		}
		return NewMismatchError("Bytes", b, state.Position)
	}
}

// NotBytes matches any byte other than the given bytes.
func NotBytes(b ...byte) Parser {
	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return NewTraceError("NotBytes", err)
		}
		for i := range b {
			if state.Buffer[state.Index] == b[i] {
				return NewNotMismatchError("NotBytes", b, state.Position)
			}
		}
		result.Value = state.Buffer[state.Index]
		state.Advance(1)
		return nil
	}
}

// ByteRange matches a byte between a given range.
func ByteRange(begin, end byte) Parser {
	if begin > end {
		panic(fmt.Errorf("byte `0x%x` is greater than `0x%x`", begin, end))
	}
	if begin == end {
		return Byte(begin)
	}
	expected := []byte{begin, '-', end}
	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return NewTraceError("ByteRange", err)
		}
		b := state.Buffer[state.Index]
		if b <= begin || end <= b {
			return NewMismatchError("ByteRange", expected, state.Position)
		}
		result.Value = b
		state.Advance(1)
		return nil
	}
}

// NotByteRange matches a byte outside a given range.
func NotByteRange(begin, end byte) Parser {
	if begin > end {
		panic(fmt.Errorf("byte `0x%x` is greater than `0x%x`", begin, end))
	}
	if begin == end {
		return Byte(begin)
	}
	expected := []byte{begin, '-', end}
	return func(state *State, result *Result) error {
		if err := state.Want(1); err != nil {
			return NewTraceError("NotByteRange", err)
		}
		b := state.Buffer[state.Index]
		if begin <= b && b <= end {
			return NewNotMismatchError("NotByteRange", expected, state.Position)
		}
		result.Value = b
		state.Advance(1)
		return nil
	}
}

// ByteSlice matches a slice of bytes.
func ByteSlice(p []byte) Parser {
	n := len(p)
	if n == 0 {
		return Epsilon
	}
	if n == 1 {
		return Byte(p[0])
	}
	return func(state *State, result *Result) error {
		if err := state.Want(n); err != nil {
			return NewTraceError("ByteSlice", err)
		}
		if !matchByteSlice(p, state.Buffer[state.Index:]) {
			return NewMismatchError("ByteSlice", p, state.Position)
		}
		result.Value = p
		state.Advance(n)
		return nil
	}
}

// Rune matches a given rune.
func Rune(r rune) Parser {
	n := utf8.RuneLen(r)
	p := make([]byte, n)
	utf8.EncodeRune(p, r)
	return func(state *State, result *Result) error {
		if err := state.Want(n); err != nil {
			return NewTraceError("Rune", err)
		}
		if !matchByteSlice(p, state.Buffer[state.Index:]) {
			return NewMismatchError("Rune", p, state.Position)
		}
		result.Value = r
		state.Advance(n)
		return nil
	}
}

// NotRune matches any rune other than the given rune.
func NotRune(r rune) Parser {
	n := utf8.RuneLen(r)
	p := make([]byte, n)
	utf8.EncodeRune(p, r)
	return func(state *State, result *Result) error {
		if err := state.Want(n); err != nil {
			return NewTraceError("NotRune", err)
		}
		if matchByteSlice(p, state.Buffer[state.Index:]) {
			return NewNotMismatchError("NotRune", p, state.Position)
		}
		state.Want(utf8.MaxRune)
		v, size := utf8.DecodeRune(state.Buffer[state.Index:])
		if v == utf8.RuneError {
			return NewParserError("failed to decode rune", state.Position)
		}
		result.Value = v
		state.Advance(size)
		return nil
	}
}

// Runes matches one of the given runes.
func Runes(r ...rune) Parser {
	n := len(r)
	p := make([][]byte, n)
	for i := range r {
		m := utf8.RuneLen(r[i])
		p[i] = make([]byte, m)
		utf8.EncodeRune(p[i], r[i])
	}
	return func(state *State, result *Result) error {
		for i := range p {
			if err := state.Want(len(p[i])); err != nil {
				return NewTraceError("Runes", err)
			}
			if matchByteSlice(p[i], state.Buffer[state.Index:]) {
				result.Value = r[i]
				state.Advance(len(p[i]))
				return nil
			}
		}
		return NewMismatchError("Runes", []byte(string(r)), state.Position)
	}
}

// NotRunes matches any rune other than the given runes.
func NotRunes(r ...rune) Parser {
	n := len(r)
	p := make([][]byte, n)
	for i := range r {
		m := utf8.RuneLen(r[i])
		p[i] = make([]byte, m)
		utf8.EncodeRune(p[i], r[i])
	}
	return func(state *State, result *Result) error {
		for i := range p {
			if err := state.Want(len(p[i])); err != nil {
				return NewTraceError("NotRunes", err)
			}
			if matchByteSlice(p[i], state.Buffer[state.Index:]) {
				return NewNotMismatchError("NotRunes", []byte(string(r)), state.Position)
			}
		}
		state.Want(utf8.MaxRune)
		v, size := utf8.DecodeRune(state.Buffer[state.Index:])
		if v == utf8.RuneError {
			return NewParserError("failed to decode rune", state.Position)
		}
		result.Value = v
		state.Advance(size)
		return nil
	}
}

// RuneRange matches a range of runes.
func RuneRange(begin, end rune) Parser {
	if begin > end {
		panic(fmt.Errorf("rune `0x%x` is greater than `0x%x`", begin, end))
	}
	if begin == end {
		return Rune(begin)
	}
	bb := make([]byte, utf8.RuneLen(begin))
	be := make([]byte, utf8.RuneLen(end))
	utf8.EncodeRune(bb, begin)
	utf8.EncodeRune(be, end)
	expected := append(append(bb, '-'), be...)
	return func(state *State, result *Result) error {
		state.Want(utf8.MaxRune)
		v, size := utf8.DecodeRune(state.Buffer[state.Index:])
		if v == utf8.RuneError {
			return NewParserError("failed to decode rune", state.Position)
		}
		if v < begin || end < v {
			return NewMismatchError("RuneRange", expected, state.Position)
		}
		result.Value = v
		state.Advance(size)
		return nil
	}
}

// NotRuneRange matches a range of runes.
func NotRuneRange(begin, end rune) Parser {
	if begin > end {
		panic(fmt.Errorf("rune `0x%x` is greater than `0x%x`", begin, end))
	}
	if begin == end {
		return Rune(begin)
	}
	bb := make([]byte, utf8.RuneLen(begin))
	be := make([]byte, utf8.RuneLen(end))
	utf8.EncodeRune(bb, begin)
	utf8.EncodeRune(be, end)
	expected := append(append(bb, '-'), be...)
	return func(state *State, result *Result) error {
		state.Want(utf8.MaxRune)
		v, size := utf8.DecodeRune(state.Buffer[state.Index:])
		if v == utf8.RuneError {
			return NewParserError("failed to decode rune", state.Position)
		}
		if begin <= v && v <= end {
			return NewNotMismatchError("NotRuneRange", expected, state.Position)
		}
		result.Value = v
		state.Advance(size)
		return nil
	}
}

// RuneSlice matches a slice of bytes.
func RuneSlice(r []rune) Parser {
	n := len(r)
	if n == 0 {
		return Epsilon
	}
	if n == 1 {
		return Rune(r[0])
	}
	p := make([]byte, 0)
	for i := range r {
		b := make([]byte, utf8.RuneLen(r[i]))
		utf8.EncodeRune(b, r[i])
		p = append(p, b...)
	}
	return func(state *State, result *Result) error {
		if err := state.Want(n); err != nil {
			return NewTraceError("RuneSlice", err)
		}
		for i := range p {
			if p[i] != state.Buffer[state.Index+i] {
				return NewMismatchError("RuneSlice", p, state.Position)
			}
		}
		result.Value = r
		state.Advance(n)
		return nil
	}
}

// String matches a given string.
func String(s string) Parser {
	bs := []byte(s)
	n := len(bs)
	return func(state *State, result *Result) error {
		if err := state.Want(n); err != nil {
			return NewTraceError("String", err)
		}
		for i := range bs {
			if bs[i] != state.Buffer[state.Index+i] {
				return NewMismatchError("String", bs, state.Position)
			}
		}
		result.Value = s
		state.Advance(n)
		return nil
	}
}

// Not transforms the given ParserLike to its Not* counterpart.
func Not(q ParserLike) Parser {
	switch p := q.(type) {
	case byte:
		return NotByte(p)
	case rune:
		return NotRune(p)
	default:
		panic(fmt.Errorf("cannot infer the Not* counterpart for type %T", p))
	}
}
