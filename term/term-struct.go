package term // import "github.com/amalog/go/term"

import "fmt"

type Struct struct {
	Context Var
	Name    Atom
	Args    Seq
	Data    Database
}

func (s *Struct) String() string {
	return fmt.Sprintf("%#v;\n", s)
}
