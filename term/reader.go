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
	x, err := r.s.Scan()
	if err != nil {
		return nil, err
	}

	switch x.Class {
	case scanner.Atom:
		y, err := r.s.Scan()
		if err != nil {
			return nil, err
		}

		switch y.Class {
		case scanner.Punct:
			if y.Text == ";" {
				return NewAtom(x.Text), nil
			}
		}
		return nil, &ErrUnexpectedToken{y}
	case scanner.String:
		if len(x.Text) < 2 {
			return nil, &Err{x, "string token too short"}
		}
		if x.Text[0] != '"' {
			return nil, &Err{x, "string missing opening quote"}
		}
		if x.Text[len(x.Text)-1] != '"' {
			return nil, &Err{x, "string missing closing quote"}
		}
		text := x.Text[1 : len(x.Text)-1]

		y, err := r.s.Scan()
		if err != nil {
			return nil, err
		}

		switch y.Class {
		case scanner.Punct:
			if y.Text == ";" {
				return NewString(text), nil
			}
		}
		return nil, &ErrUnexpectedToken{y}
	case scanner.Var:
		y, err := r.s.Scan()
		if err != nil {
			return nil, err
		}

		switch y.Class {
		case scanner.Punct:
			if y.Text == ";" {
				term := &Var{
					Name:  x.Text,
					Value: nil,
				}
				return term, nil
			}
		}
		return nil, &ErrUnexpectedToken{y}
	default:
		return nil, &ErrUnexpectedToken{x}
	}

	return nil, nil
}
