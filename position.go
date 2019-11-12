package pars

import "fmt"

// Position represents the line and byte numbers.
type Position struct {
	Line int
	Byte int
}

// Head tests if the position is at the head.
func (p Position) Head() bool {
	return p.Line == 0 && p.Byte == 0
}

// String returns a formatted position.
func (p Position) String() string {
	return fmt.Sprintf("line %d byte %d", p.Line+1, p.Byte+1)
}
