package scanner // import "github.com/amalog/go/scanner"
import "fmt"

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
