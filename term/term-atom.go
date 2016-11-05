package term // import "github.com/amalog/go/term"
import "fmt"

type Atom string

func NewAtom(s string) Term {
	return Atom(s)
}

func (a Atom) String() string {
	return fmt.Sprintf("%s;\n", string(a))
}
