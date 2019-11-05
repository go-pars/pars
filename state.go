package pars

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

func notEOF(err error) bool { return err != nil && err != io.EOF }

const (
	stackGrowthSize = 16
	bufferReadSize  = 4096
)

// Position represents the line and byte numbers.
type Position struct {
	Line int
	Byte int
}

func (p Position) Head() bool {
	return p.Line == 0 && p.Byte == 0
}

func (p Position) String() string {
	return fmt.Sprintf("line %d byte %d", p.Line+1, p.Byte+1)
}

type frame struct {
	Index    int
	Position Position
}

type stack struct {
	frames []frame
	index  int
}

func newStack() *stack { return &stack{make([]frame, stackGrowthSize), 0} }

func (s stack) Len() int { return s.index }

func (s *stack) Push(index int, position Position) {
	if len(s.frames) == s.index {
		s.frames = append(s.frames, make([]frame, stackGrowthSize)...)
	}
	s.frames[s.index] = frame{index, position}
	s.index++
}

func (s *stack) Pop() (int, Position) {
	if s.index == 0 {
		panic("Pop called on empty stack")
	}
	s.index--
	f := s.frames[s.index]
	return f.Index, f.Position
}

func (s *stack) Clear() {
	s.index = 0
}

// State represents the parser state.
type State struct {
	reader   io.Reader
	buffer   []byte
	index    int
	isEOF    bool
	wanted   int
	position Position
	frames   *stack
}

// NewState creates a new state from the given io.Reader.
func NewState(r io.Reader) *State {
	return &State{
		reader:   r,
		buffer:   make([]byte, 0),
		index:    0,
		isEOF:    false,
		wanted:   0,
		position: Position{0, 0},
		frames:   newStack(),
	}
}

// FromString creates a new state from the given string.
func FromString(s string) *State { return NewState(strings.NewReader(s)) }

// FromBytes creates a new state from the given bytes.
func FromBytes(p []byte) *State { return NewState(bytes.NewBuffer(p)) }

// Read bytes from the state.
func (s *State) Read(p []byte) (int, error) {
	l := len(p)

	// Attempt to read full length of p.
	if err := s.Want(l); notEOF(err) {
		return 0, err
	}

	// Check for the number of bytes left in the buffer.
	n := len(s.buffer)
	if n < l {
		l = n
	}

	copy(p, s.buffer)
	s.Advance()
	return n, nil
}

func (s *State) fill() (int, error) {
	p := make([]byte, bufferReadSize)
	n, err := s.reader.Read(p)
	if n > 0 || err == nil {
		s.buffer = append(s.buffer, p[:n]...)
	}
	return n, err
}

// Want queries the state for the number of bytes wanted.
// If there are not enough bytes in the buffer, additional bytes are read
// from the io.Reader.
func (s *State) Want(n int) error {
	// There are not enough bytes left in the buffer.
	if s.index+n > len(s.buffer) {
		// The io.Reader already reached EOF.
		if s.isEOF {
			return io.EOF
		}

		// Read the next block of bytes.
		m, err := s.fill()
		if err == io.EOF {
			s.isEOF = true
		} else if err != nil {
			return err
		}

		// Still not enough bytes.
		if m < n {
			s.wanted = m
			return io.EOF
		}
	}

	s.wanted = n
	return nil
}

// Head returns the first byte in the buffer.
func (s State) Head() byte { return s.buffer[s.index] }

// Is checks if the head is the given byte.
func (s State) Is(c byte) bool { return s.buffer[s.index] == c }

// Buffer returns the byte slice which is accessible to the user.
func (s State) Buffer() []byte { return s.buffer[s.index : s.index+s.wanted] }

// Dump returns the entire remaining buffer content.
func (s State) Dump() []byte { return s.buffer[s.index:] }

// Advance the index by the amount given in a previous Want call.
func (s *State) Advance() {
	n := s.wanted
	if n == 0 {
		panic("no previous call to Want")
	}
	for _, b := range s.buffer[s.index : s.index+n] {
		if b == '\n' {
			s.position.Line++
			s.position.Byte = 0
		} else {
			s.position.Byte++
		}
	}
	s.index += n
	s.wanted = 0
	s.autoclear()
}

// Position returns the current position of the state.
func (s State) Position() Position { return s.position }

// Push the current state frame into the internal stack.
func (s *State) Push() { s.frames.Push(s.index, s.position) }

// Pop the latest frame from the internal stack.
func (s *State) Pop() error {
	if s.frames.Len() == 0 {
		return errors.New("state stack empty")
	}
	s.index, s.position = s.frames.Pop()
	s.autoclear()
	return nil
}

// Drop the latest frame from the internal stack.
func (s *State) Drop() { s.frames.Pop(); s.autoclear() }

func (s *State) autoclear() {
	if s.frames.Len() == 0 {
		s.Clear()
	}
}

// Clear the state buffer and index.
func (s *State) Clear() {
	s.buffer = s.buffer[s.index:]
	s.index = 0
	s.frames.Clear()
}
