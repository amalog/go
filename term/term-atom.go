package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"fmt"
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
	fmt.Fprintf(w, "%s,\n", string(a))
}

func (a Atom) Name() string {
	return string(a)
}
