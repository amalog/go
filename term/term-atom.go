package term // import "github.com/amalog/go/term"
import "fmt"

type Atom string

func NewAtom(s string) Atom {
	return Atom(s)
}

func (a Atom) String() string {
	return fmt.Sprintf("%s,\n", string(a))
}

func (a Atom) Name() string {
	return string(a)
}
