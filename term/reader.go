package term // import "github.com/amalog/go/term"

import (
	"io"

	"github.com/amalog/go/scanner"
)

// central definition of the term terminator
const terminator = ","

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
		return r.readAtomOrStruct(nil, x)
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

		t := NewString(text)
		return r.terminate(t)
	case scanner.Var:
		y, err := r.s.Scan()
		if err != nil {
			return nil, err
		}

		switch y.Class {
		case scanner.Punct:
			switch y.Text {
			case terminator:
				term := &Var{
					Name:  x.Text,
					Value: nil,
				}
				return term, nil
			case ".":
				z, err := r.s.Scan()
				if err != nil {
					return nil, err
				}
				if z.Class == scanner.Atom {
					return r.readAtomOrStruct(x, z)
				}
			}
		}
		return nil, &ErrUnexpectedToken{y}
	default:
		return nil, &ErrUnexpectedToken{x}
	}

	return nil, nil
}

func (r *Reader) readAtomOrStruct(context, name *scanner.Token) (Term, error) {
	y, err := r.s.Scan()
	if err != nil {
		return nil, err
	}

	switch y.Class {
	case scanner.Punct:
		switch y.Text {
		case terminator:
			return NewAtom(name.Text), nil
		case "(":
			args, err := r.readSeq(y.Text) // consumes closing ')'
			if err != nil {
				return nil, err
			}

			z, err := r.s.Scan()
			if err != nil {
				return nil, err
			}

			t := &Struct{
				Name: NewAtom(name.Text),
				Args: NewSeq(args),
			}
			if context != nil {
				t.Context = Var{Name: context.Text}
			}

			switch z.Class {
			case scanner.Punct:
				switch z.Text {
				case terminator:
					return t, nil
				case "{":
					data, err := r.readSeq(z.Text) // consumes closing '}'
					if err != nil {
						return nil, err
					}
					t.Data = NewDb(data)
					return r.terminate(t)
				}
			}
			return nil, &ErrUnexpectedToken{z}
		}
	}
	return nil, &ErrUnexpectedToken{y}
}

// reads a sequence of terms. should be called immediately after the '('
// (or '}') token is consumed.  consumes the closing ')' (or '}') token.
func (r *Reader) readSeq(open string) ([]Term, error) {
	args := make([]Term, 0)

	var close string
	switch open {
	case "(":
		close = ")"
	case "{":
		close = "}"
	}

ARGS:
	for {
		t, err := r.Read()
		switch e := err.(type) {
		case nil:
			args = append(args, t)
		case *ErrUnexpectedToken:
			if e.Token.Class == scanner.Punct && e.Token.Text == close {
				break ARGS
			}
			return nil, err
		default:
			return nil, err
		}
	}

	return args, nil
}

// consume a terminator. if successful return the given term
func (r *Reader) terminate(t Term) (Term, error) {
	x, err := r.s.Scan()
	if err == io.EOF {
		return nil, &ErrUnexpectedEof{r.s.Pos()}
	}
	if err != nil {
		return nil, err
	}

	if x.Class == scanner.Punct && x.Text == terminator {
		return t, nil
	}
	return nil, &ErrUnexpectedToken{x}
}
