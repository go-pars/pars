package pars

import (
	"bufio"
	"io"

	ascii "gopkg.in/ktnyt/ascii.v1"
)

// Reader is a special io.Reader that will skip all whitespaces unless it is
// a part of a string literal (quoted by a ", ', or `).
type Reader struct {
	reader  *bufio.Reader
	quoted  bool
	escaped bool
}

// NewReader creates a new reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{bufio.NewReader(r), false, false}
}

// Read satisfies the io.Reader interface.
func (r *Reader) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		c, err := r.reader.ReadByte()

		if err != nil {
			return n, err
		}

		if !r.escaped && ascii.IsQuote(c) {
			r.quoted = !r.quoted
		}

		if r.quoted || !ascii.IsSpace(c) {
			p[n] = c
			n++
		}

		r.escaped = c == '\\'
	}

	return n, nil
}
