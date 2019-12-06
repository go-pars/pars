package pars

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

const (
	bufferReadSize = 4096
)

// State represents a parser state, which is a convenience wrapper for an
// io.Reader object with buffering and backtracking.
type State struct {
	rd  io.Reader
	buf []byte
	off int
	err error
	req int
	pos Position
	stk *stack
}

// NewState creates a new state from the given io.Reader.
func NewState(r io.Reader) *State {
	return &State{
		rd:  r,
		buf: make([]byte, 0),
		off: 0,
		err: nil,
		req: -1,
		pos: Position{0, 0},
		stk: newStack(),
	}
}

// FromString creates a new state from the given string.
func FromString(s string) *State { return NewState(strings.NewReader(s)) }

// FromBytes creates a new state from the given bytes.
func FromBytes(p []byte) *State { return NewState(bytes.NewBuffer(p)) }

// Read satisfies the io.Reader interface.
func (s *State) Read(p []byte) (int, error) {
	err := s.Request(len(p))
	n := copy(p, s.buf)
	s.Advance()
	return n, err
}

// Request checks if the state contains at least the given number of bytes,
// additionally reading from the io.Reader object as necessary when the
// internal buffer is exhausted. If the call to Read for the io.Reader object
// returns an error, Request will return the corresponding error.
func (s *State) Request(n int) error {
	// There are not enough bytes left in the buffer.
	for len(s.buf) < s.off+n {
		// The io.Reader object previously returned an error.
		if s.err != nil {
			return s.err
		}

		// Read the next block of bytes.
		p := make([]byte, len(s.buf)+bufferReadSize)
		l := copy(p, s.buf)
		m, err := s.rd.Read(p[l:])
		s.buf = p[:l+m]
		if err != nil {
			s.err = err
			// Still not enough bytes.
			if len(s.buf) < s.off+n {
				return s.err
			}
		}
	}

	s.req = n
	return nil
}

// Advance the state by the amount given in a previous Request call.
func (s *State) Advance() {
	if s.req == -1 {
		panic("no previous call to Request")
	}
	for _, b := range s.buf[s.off : s.off+s.req] {
		if b == '\n' {
			s.pos.Line++
			s.pos.Byte = 0
		} else {
			s.pos.Byte++
		}
	}
	s.off += s.req
	s.req = -1
	s.autoclear()
}

// Buffer returns the range of bytes guaranteed by a Request call.
func (s State) Buffer() []byte { return s.buf[s.off : s.off+s.req] }

// Dump returns the entire remaining buffer content. Note that the returned
// byte slice will not always contain the entirety of the bytes that can be
// read by the io.Reader object.
func (s State) Dump() []byte { return s.buf[s.off:] }

// Offset returns the current state offset.
func (s State) Offset() int { return s.off }

// Position returns the current line and byte position of the state.
func (s State) Position() Position { return s.pos }

// Push the current state position for backtracking.
func (s *State) Push() { s.stk.Push(s.off, s.pos) }

// Pushed tests if the state has been pushed at least once.
func (s State) Pushed() bool { return !s.stk.Empty() }

// Pop will backtrack to the most recently pushed state.
func (s *State) Pop() {
	if !s.stk.Empty() {
		s.off, s.pos = s.stk.Pop()
		s.autoclear()
	}
}

// Drop will discard the most recently pushed state.
func (s *State) Drop() {
	if !s.stk.Empty() {
		s.stk.Pop()
		s.autoclear()
	}
}

func (s *State) autoclear() {
	if s.stk.Empty() {
		s.Clear()
	}
}

// Clear will discard the buffer contents prior to the current state offset
// and drop all pushed states.
func (s *State) Clear() {
	s.buf = s.buf[s.off:]
	s.off = 0
	s.stk.Reset()
}

// Skip the given state for the given number of bytes.
func Skip(state *State, n int) error {
	if err := state.Request(n); err != nil {
		return err
	}
	state.Advance()
	return nil
}

// Next attempts to retrieve the next byte in the given state.
func Next(state *State) (byte, error) {
	if err := state.Request(1); err != nil {
		return 0, err
	}
	return state.Buffer()[0], nil
}

// Trail will return the extent of the state buffer spanning from the most
// recently pushed state position to the current state position.
func Trail(state *State) ([]byte, error) {
	if !state.Pushed() {
		return nil, errors.New("failed to backtrack")
	}
	off := state.Offset()
	state.Pop()
	n := off - state.Offset()
	if state.Request(n) != nil {
		panic("logical error: Request failed")
	}
	p := state.Buffer()
	state.Advance()
	return p, nil
}
