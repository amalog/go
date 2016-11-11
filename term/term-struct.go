package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"io"
)

type Struct struct {
	Context Var
	Name    Atom
	Args    Seq
	Data    Database
}

func (s *Struct) String() string {
	buf := new(bytes.Buffer)
	s.Format(buf, Style{})
	return buf.String()
}

func (s *Struct) Format(w io.Writer, style Style) {
	io.WriteString(w, s.Name.Name())
	io.WriteString(w, "(")
	s.Args.Format(w, style)
	io.WriteString(w, ")")

	if len(s.Data) > 0 {
		panic("can't yet display structs with a database")
	}

	if !style.OmitTerminator {
		io.WriteString(w, ",")
	}
	if !style.OmitNewline {
		io.WriteString(w, "\n")
	}
}
