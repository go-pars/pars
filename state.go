package pars

import (
	"io"
	"strings"
)

const (
	stackGrowthSize = 16
	bufferReadSize  = 4096
)

type stack struct {
	buffer []int
	index  int
}

func newStack() *stack {
	return &stack{buffer: make([]int, stackGrowthSize)}
}

func (s *stack) Push(n int) {
	if len(s.buffer) == s.index {
		s.buffer = append(s.buffer, make([]int, stackGrowthSize)...)
	}
	s.buffer[s.index] = n
	s.index++
}

func (s *stack) Pop() (int, bool) {
	if s.index == 0 {
		return 0, false
	}
	s.index--
	return s.buffer[s.index], true
}

func (s *stack) Clear() {
	s.buffer = make([]int, stackGrowthSize)
	s.index = 0
}

// State represents the parser state.
type State struct {
	reader io.Reader
	marks  *stack
	isEOF  bool
	isDry  bool

	Buffer   []byte
	Index    int
	Position int
}

// NewState returns a new State for io.Reader.
func NewState(r io.Reader) *State {
	s := &State{reader: r, marks: newStack()}
	if err := s.fill(); err != nil {
		if err != io.EOF {
			panic("could not read from io.Reader on first read")
		}
	}
	return s
}

// FromString returns a new State for given string.
func FromString(s string) *State {
	return NewState(strings.NewReader(s))
}

// Read bytes from the state.
func (s *State) Read(p []byte) (int, error) {
	l := len(p)
	err := s.Want(l)
	if err != nil && err != io.EOF {
		return 0, err
	}
	copy(p, s.Buffer[s.Index:])
	n := len(s.Buffer[s.Index:])
	if l < n {
		n = l
	}
	s.Advance(n)
	return n, err
}

func (s *State) fill() error {
	next := make([]byte, bufferReadSize)
	n, err := s.reader.Read(next)
	if n == 0 && err != nil {
		return err
	}
	if err == io.EOF {
		s.isEOF = true
	}
	s.Buffer = append(s.Buffer, next[:n]...)
	return nil
}

// Want tells the State how many bytes are wanted.
// Data is read from the io.Reader if there are not enough bytes.
func (s *State) Want(n int) error {
	if len(s.Buffer) < s.Index+n {
		if s.isEOF {
			return io.EOF
		}
		return s.fill()
	}
	return nil
}

// Advance the Index and Position of the state.
func (s *State) Advance(n int) {
	s.Index += n
	s.Position += n
}

// Mark the current Index for future Jump.
func (s *State) Mark() {
	s.marks.Push(s.Index)
}

// Unmark the latest Mark.
func (s *State) Unmark() bool {
	_, ok := s.marks.Pop()
	return ok
}

// Remark the latest Mark.
func (s *State) Remark() {
	s.Unmark()
	s.Mark()
}

// Jump to the latest Mark.
func (s *State) Jump() bool {
	n, ok := s.marks.Pop()
	if !ok {
		return false
	}
	s.Position -= s.Index - n
	s.Index = n
	return true
}

// Dry makes the state dry.
func (s *State) Dry() {
	s.isDry = true
}

// Wet makes the state wet.
func (s *State) Wet() {
	s.isDry = false
}

// Clear the state Buffer, Index, and Marks.
func (s *State) Clear() {
	if !s.isDry {
		s.Buffer = s.Buffer[s.Index:]
		s.Index = 0
		s.marks.Clear()
	}
}

// EOF returns true if the state reached the end of the reader.
func (s *State) EOF() bool {
	return s.isEOF
}
