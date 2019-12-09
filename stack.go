package pars

const stackGrowthSize = 16

type frame struct {
	Off int
	Pos Position
}

type stack struct {
	v []frame
	i int
}

func newStack() *stack { return &stack{make([]frame, stackGrowthSize), 0} }

func (s stack) Empty() bool { return s.i == 0 }

func (s *stack) Push(i int, position Position) {
	if s.i == len(s.v) {
		s.v = append(s.v, make([]frame, stackGrowthSize)...)
	}
	s.v[s.i] = frame{i, position}
	s.i++
}

func (s *stack) Pop() (int, Position) {
	s.i--
	f := s.v[s.i]
	return f.Off, f.Pos
}

func (s *stack) Reset() { s.i = 0 }
