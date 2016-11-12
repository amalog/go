package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"io"
)

type Atom string

func NewAtom(s string) Atom {
	return Atom(s)
}

func (a Atom) String() string {
	buf := new(bytes.Buffer)
	a.Format(buf, Style{})
	return buf.String()
}

func (a Atom) Format(w io.Writer, style Style) {
	io.WriteString(w, string(a))
	style.Terminate(w)
}

func (a Atom) Name() string {
	return string(a)
}
