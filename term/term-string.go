package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"fmt"
	"io"
)

type String string

func NewString(s string) String {
	return String(s)
}

func (s String) String() string {
	buf := new(bytes.Buffer)
	s.Format(buf, Style{})
	return buf.String()
}

func (s String) Format(w io.Writer, style Style) {
	fmt.Fprintf(w, "\"%s\",\n", string(s))
}
