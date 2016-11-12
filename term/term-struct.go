package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"io"
)

type Struct struct {
	Context Var
	Name    Atom
	Args    Seq
	Data    Db
}

func (s *Struct) String() string {
	buf := new(bytes.Buffer)
	s.Format(buf, Style{})
	return buf.String()
}

func (s *Struct) Format(w io.Writer, style Style) {
	if s.Context.Name != "" {
		io.WriteString(w, s.Context.Name)
		io.WriteString(w, ".")
	}

	io.WriteString(w, s.Name.Name())
	io.WriteString(w, "(")
	s.Args.Format(w, style)
	io.WriteString(w, ")")

	if len(s.Data) > 0 {
		io.WriteString(w, " {\n")
		style.Indent++
		s.Data.Format(w, style)
		style.Indent--
		io.WriteString(w, "}")
	}

	if !style.OmitTerminator {
		io.WriteString(w, ",")
	}
	if !style.OmitNewline {
		io.WriteString(w, "\n")
	}
}
