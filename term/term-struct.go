package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"fmt"
)

type Struct struct {
	Context Var
	Name    Atom
	Args    Seq
	Data    Database
}

func (s *Struct) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s(%s)", s.Name.Name(), s.Args)
	if len(s.Data) > 0 {
		panic("can't yet display structs with a database")
	}
	buf.WriteString(",")
	return buf.String()
}
