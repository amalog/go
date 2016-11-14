package scanner // import "github.com/amalog/go/scanner"

import (
	"fmt"
	"io"
)

type Class int

const (
	Atom  Class = iota
	Eof   Class = iota
	Num   Class = iota
	Punct Class = iota
	Var   Class = iota

	String Class = iota

	// classes used internally
	nl Class = iota
)

const eof rune = -1

type Position struct {
	Filename string
	Line     int
	Column   int
}

func (pos Position) String() string {
	f := pos.Filename
	if f == "" {
		f = "<input>"
	}
	return fmt.Sprintf("%s:%d:%d", f, pos.Line, pos.Column)
}

type SyntaxError struct {
	Position Position
	Message  string
}

func (err *SyntaxError) Error() string {
	return fmt.Sprintf("%s: %s", err.Position, err.Message)
}

type Token struct {
	Class    Class
	Position Position
	Text     string
}

func (t *Token) String() string {
	switch t.Class {
	case Atom:
		return fmt.Sprintf("atom(%s)", t.Text)
	case Eof:
		return "eof"
	case Num:
		return fmt.Sprintf("num(%s)", t.Text)
	case Punct:
		return fmt.Sprintf("punct(%s)", t.Text)
	case String:
		return fmt.Sprintf("string(%s)", t.Text)
	case Var:
		return fmt.Sprintf("var(%s)", t.Text)
	default:
		return "<unknown token class>"
	}
}

type Scanner struct {
	r io.RuneScanner

	filename string
	err      *SyntaxError

	line       int
	column     int
	prevLine   int // before calling next()
	prevColumn int // before calling next()

	prevToken *Token
}

func New(r io.RuneScanner) *Scanner {
	s := &Scanner{
		r:      r,
		line:   1,
		column: 0,
	}
	return s
}

// Scan returns the next token in the input.  Returns an infinite stream of Eof
// tokens once reaching the end of input.
func (s *Scanner) Scan() (*Token, error) {
	t, err := s.scan()
	if err != nil {
		s.prevToken = nil
		return nil, err
	}

	// handle comma insertion
	s.insertComma(t)
	if t.Class == nl { // newlines are for internal use only. skip them
		return s.Scan()
	}

	s.prevToken = t
	return t, err
}

func (s *Scanner) scan() (*Token, error) {
	var ch rune
	for {
		ch = s.next()
		if ch == eof {
			if s.err == nil {
				return &Token{Class: Eof, Position: s.Pos()}, nil
			}
			return nil, s.err
		}

		if ch >= 'a' && ch <= 'z' { // atom
			return s.scanAtom(ch)
		} else if ch >= 'A' && ch <= 'Z' { // variable
			return s.scanVariable(ch)
		} else if ch >= '0' && ch <= '9' { // number
			return s.scanNumber(ch)
		} else if ch == '"' { // string
			return s.scanString(ch)
		} else if ch == ' ' { // space
			s.skipSpace()
			continue
		} else if ch == '\n' { // newline
			return &Token{Class: nl, Position: s.Pos()}, nil
		} else {
			break
		}
	}

	t := &Token{
		Class:    Punct,
		Position: s.Pos(),
		Text:     string([]rune{ch}),
	}
	return t, nil
}

func (s *Scanner) insertComma(t *Token) {
	// don't insert comma as first token
	if s.prevToken == nil {
		return
	}

	// certain punctuation prohibits comma insertion
	if s.prevToken.Class == Punct {
		switch s.prevToken.Text {
		case ",", "(", "{":
			return
		}
	}

	switch t.Class {
	case Eof:
		if s.prevToken.Class != Eof { // before first EOF token
			t.Class = Punct
			t.Text = ","
		}
	case nl:
		t.Text = ","
		t.Class = Punct
	case Punct:
		if t.Text == ")" || t.Text == "}" { // before closing a seq or db
			t.Text = ","
			s.back()
		}
	}
}

func (s *Scanner) peek() rune {
	ch := s.next()
	if ch == eof {
		return ch
	}
	s.back()
	return ch
}

func (s *Scanner) back() {
	err := s.r.UnreadRune()
	if err != nil {
		panic(err)
	}

	s.line = s.prevLine
	s.column = s.prevColumn
}

// consumes the next character in the stream. stores errors in s.err
func (s *Scanner) next() rune {
	s.prevLine = s.line
	s.prevColumn = s.column

	ch, _, err := s.r.ReadRune()
	if err == io.EOF {
		return eof
	}
	if err != nil {
		panic(err)
	}

	// handle prohibited characters
	if ch == '\t' || ch == '\r' {
		s.err = s.prohibitedCharacter(ch)
		return eof
	}

	// update position information
	if ch == '\n' {
		s.line++
		s.column = 0
	} else {
		s.column++
	}

	return ch
}

func (s *Scanner) prohibitedCharacter(ch rune) *SyntaxError {
	var name string
	switch ch {
	case '\t':
		name = "tab"
	case '\r':
		name = "carriage return"
	}
	return &SyntaxError{
		Position: s.Pos(),
		Message:  fmt.Sprintf("The %s character is prohibited", name),
	}
}

func (s *Scanner) Pos() Position {
	return Position{
		Filename: s.filename,
		Line:     s.line,
		Column:   s.column,
	}
}

func (s *Scanner) skipSpace() {
	for {
		ch := s.next()
		if ch == ' ' {
			continue
		}
		if ch != eof {
			s.back()
		}
		break
	}
}

func (s *Scanner) scanAtom(ch rune) (*Token, error) {
	chars := make([]rune, 0)

	pos := s.Pos()
CH:
	for {
		switch {
		case ch >= 'a' && ch <= 'z', ch == '_':
			chars = append(chars, ch)
		case ch == eof:
			break CH
		default:
			s.back()
			break CH
		}

		ch = s.next()
	}

	t := &Token{
		Class:    Atom,
		Position: pos,
		Text:     string(chars),
	}
	return t, nil
}

func (s *Scanner) scanVariable(ch rune) (*Token, error) {
	chars := make([]rune, 0)

	pos := s.Pos()
CH:
	for {
		switch {
		case ch >= 'A' && ch <= 'Z':
			if len(chars) > 0 {
				prev := chars[len(chars)-1]
				if prev < 'a' || prev > 'z' {
					err := &SyntaxError{
						Position: s.Pos(),
						Message:  fmt.Sprintf("variable names may not have consecutive uppercase letters, got %c", ch),
					}
					return nil, err
				}
			}
			chars = append(chars, ch)
		case ch >= 'a' && ch <= 'z':
			if len(chars) == 0 {
				panic("called scanVariable without upper case letter next in stream")
			}
			chars = append(chars, ch)
		case ch == eof:
			break CH
		default:
			s.back()
			break CH
		}

		ch = s.next()
	}

	t := &Token{
		Class:    Var,
		Position: pos,
		Text:     string(chars),
	}
	return t, nil
}

func (s *Scanner) scanNumber(ch rune) (*Token, error) {
	x, err := s.scanInteger(ch)
	if err != nil {
		return nil, err
	}

	switch ch := s.next(); ch {
	case '.':
		ch = s.next()
		y, err := s.scanInteger(ch)
		if err != nil {
			return nil, err
		}
		t := &Token{
			Class:    Num,
			Position: x.Position,
			Text:     x.Text + "." + y.Text,
		}
		return t, nil
	default:
		t := &Token{
			Class:    Num,
			Position: x.Position,
			Text:     x.Text,
		}
		return t, nil
	}
}

func (s *Scanner) scanInteger(ch rune) (*Token, error) {
	chars := make([]rune, 0)

	pos := s.Pos()
CH:
	for {
		switch {
		case ch >= '0' && ch <= '9', ch == '_':
			chars = append(chars, ch)
		case ch == '(', ch == ')', ch == ',', ch == '.', ch == ' ':
			s.back()
			break CH
		case ch == eof:
			break CH
		default:
			err := &SyntaxError{
				Position: s.Pos(),
				Message:  fmt.Sprintf("Unexpected number character: '%c'", ch),
			}
			return nil, err
		}

		ch = s.next()
	}

	t := &Token{
		Class:    Num, // not really
		Position: pos,
		Text:     string(chars),
	}
	return t, nil
}

func (s *Scanner) scanString(ch rune) (*Token, error) {
	chars := make([]rune, 0)

	// consume opening quote character
	if ch != '"' {
		panic("scanString without a double quote character to start")
	}
	chars = append(chars, ch)

	ch = s.next()
	pos := s.Pos()
	for {
		chars = append(chars, ch)
		if ch == '"' {
			break
		}
		ch = s.next()
	}

	t := &Token{
		Class:    String,
		Position: pos,
		Text:     string(chars),
	}
	return t, nil
}
