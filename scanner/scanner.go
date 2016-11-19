package scanner // import "github.com/amalog/go/scanner"

import (
	"fmt"
	"io"
)

const eof rune = -1

type Scanner struct {
	r io.RuneScanner

	err *SyntaxError

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
		} else if ch == '#' { // end of line comment
			return s.scanEolComment(ch)
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

	// certain tokens prohibit comma insertion
	switch s.prevToken.Class {
	case Punct:
		switch s.prevToken.Text {
		case ",", "(", "{":
			return
		}
	case Comment:
		return
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
		Line:   s.line,
		Column: s.column,
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
		case ch >= 'A' && ch <= 'Z', ch >= 'a' && ch <= 'z':
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
	chars := make([]rune, 0)

	pos := s.Pos()
CH:
	for {
		switch {
		case ch >= '0' && ch <= '9', ch == '_', ch == '.':
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
		Class:    Num,
		Position: pos,
		Text:     string(chars),
	}
	return t, nil
}

func (s *Scanner) scanString(ch rune) (*Token, error) {
	chars := []rune{ch}

	pos := s.Pos()
	for {
		ch = s.next()
		chars = append(chars, ch)
		if ch == '"' {
			break
		}
		if ch == eof {
			return nil, &SyntaxError{
				Position: pos,
				Message:  "Runaway string",
			}
		}
	}

	t := &Token{
		Class:    String,
		Position: pos,
		Text:     string(chars),
	}
	return t, nil
}

func (s *Scanner) scanEolComment(ch rune) (*Token, error) {
	chars := []rune{ch}

	pos := s.Pos()
	for {
		ch = s.next()
		if ch == '\n' || ch == eof {
			break
		}
		chars = append(chars, ch)
	}

	t := &Token{
		Class:    Comment,
		Position: pos,
		Text:     string(chars),
	}
	return t, nil
}
