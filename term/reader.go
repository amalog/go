package term // import "github.com/amalog/go/term"

import (
	"io"

	"github.com/amalog/go/scanner"
)

type Reader struct {
	s *scanner.Scanner
}

func NewReader(r io.RuneScanner) *Reader {
	self := new(Reader)
	self.s = scanner.New(r)
	return self
}

func (r *Reader) Read() (Term, error) {
	return nil, nil
}
