package term // import "github.com/amalog/go/term"

import "io"

type Struct struct {
	Context *Var
	Name    Atom
	Args    Seq
	Data    Db
}

func (s *Struct) Format(w io.Writer, style Style) {
	if s.Context != nil && s.Context.Name != "" {
		io.WriteString(w, s.Context.Name)
		io.WriteString(w, ".")
	}

	io.WriteString(w, s.Name.Name())

	nilary := len(s.Args) == 0 && len(s.Data) == 0
	if len(s.Args) > 0 || nilary {
		io.WriteString(w, "(")
		s.Args.Format(w, style)
		io.WriteString(w, ")")
	}

	if len(s.Data) > 0 {
		io.WriteString(w, " {\n")
		style.Indent++
		s.Data.Format(w, style)
		style.Indent--
		style.WriteIndent(w)
		io.WriteString(w, "}")
	}

	style.Terminate(w)
}
