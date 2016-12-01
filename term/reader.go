package term // import "github.com/amalog/go/term"

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/amalog/go/scanner"
)

// central definition of the term terminator
const terminator = ","

type Reader struct {
	s *scanner.Scanner
}

func NewReader(r io.Reader) *Reader {
	var rs io.RuneScanner
	if x, ok := r.(io.RuneScanner); ok {
		rs = x
	} else {
		rs = bufio.NewReader(r)
	}
	self := new(Reader)
	self.s = scanner.New(rs)
	return self
}

// ReadAllAsTerms reads all terms.
func ReadAllAsTerms(r io.Reader) ([]Term, error) {
	reader := NewReader(r)
	terms := make([]Term, 0)
	for {
		t, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}
	return terms, nil
}

// ReadAllAsTerm reads all content as a single, root term.  The term's name is
// "amalog".  All terms read from the content are in its db.
func ReadAllAsTerm(r io.Reader) (Term, error) {
	terms, err := ReadAllAsTerms(r)
	if err != nil {
		return nil, err
	}

	name, err := NewAtom("amalog")
	if err != nil {
		panic(err) // "amalog" should always be a valid atom
	}
	t := &Struct{
		Name: name,
		Data: NewDb(terms),
	}
	return t, nil
}

// ReadFileAsTerms reads all terms from a single file.
func ReadFileAsTerms(path string) ([]Term, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ReadAllAsTerms(file)
}

// ReadDirectoryAsTerm reads all files in a directory as a single, root term.
// The term's name is "amalog".  All terms read from the content are in its db.
// The order in which files are processed is undefined.
func ReadDirectoryAsTerm(dir string) (Term, error) {
	terms := make([]Term, 0)

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".ama") {
			continue
		}
		p := filepath.Join(dir, entry.Name())
		ts, err := ReadFileAsTerms(p)
		if err != nil {
			return nil, err
		}
		terms = append(terms, ts...)
	}

	name, err := NewAtom("amalog")
	if err != nil {
		panic(err) // "amalog" should always be a valid atom
	}
	t := &Struct{
		Name: name,
		Data: NewDb(terms),
	}
	return t, nil
}

// WriteAll writes a term as if it were a root term.  The term must be a
// struct.  Its name is ignored.  See ReadAll.
func WriteAll(w io.Writer, t Term) error {
	if _, ok := t.(*Struct); ok {
		style := Style{IsRoot: true}
		t.Format(w, style)
		return nil
	}
	return errors.New("WriteFile can only output a struct term")
}

func (r *Reader) Read() (Term, error) {
	x, err := r.s.Scan()
	if err != nil {
		return nil, err
	}

	switch x.Class {
	case scanner.Atom:
		return r.readAtomOrStruct(nil, x)
	case scanner.Comment:
		t, err := NewComment(x.Text)
		if err != nil {
			panic(err)
		}
		return t, nil
	case scanner.Eof:
		return nil, io.EOF
	case scanner.String:
		t := NewString(x.Text)
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
				t, err := NewVar(x.Text)
				if err != nil {
					return nil, &Err{x, err.Error()}
				}
				return t, nil
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

	var contextVar *Var
	if context != nil {
		contextVar, err = NewVar(context.Text)
		if err != nil {
			return nil, &Err{context, err.Error()}
		}
	}

	nameTerm, err := NewAtom(name.Text)
	if err != nil {
		return nil, &Err{name, err.Error()}
	}

	switch y.Class {
	case scanner.Eof:
		return nil, &ErrUnexpectedEof{r.s.Pos()}
	case scanner.Punct:
		switch y.Text {
		case terminator:
			return nameTerm, nil
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
				Context: contextVar,
				Name:    nameTerm,
				Args:    NewSeq(args),
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
		case "{":
			data, err := r.readSeq(y.Text) // consumes closing '}'
			if err != nil {
				return nil, err
			}

			t := &Struct{
				Context: contextVar,
				Name:    nameTerm,
				Data:    NewDb(data),
			}
			return r.terminate(t)
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
		if err == io.EOF {
			return nil, &ErrUnexpectedEof{r.s.Pos()}
		}
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
	if x.Class == scanner.Eof {
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
