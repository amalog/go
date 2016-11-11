package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"fmt"
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
	fmt.Fprintf(w, "%s,\n", v.Name)
}
