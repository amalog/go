package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"io"
)

type Var struct {
	Name  string
	Value *Term
}

func (v *Var) String() string {
	buf := new(bytes.Buffer)
	v.Format(buf, Style{})
	return buf.String()
}

func (v *Var) Format(w io.Writer, style Style) {
	io.WriteString(w, v.Name)

	if !style.OmitTerminator {
		io.WriteString(w, ",")
	}
	if !style.OmitNewline {
		io.WriteString(w, "\n")
	}
}
