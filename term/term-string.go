package term // import "github.com/amalog/go/term"

import (
	"bytes"
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
	io.WriteString(w, "\"")
	io.WriteString(w, string(s))
	io.WriteString(w, "\"")

	if !style.OmitTerminator {
		io.WriteString(w, ",")
	}
	if !style.OmitNewline {
		io.WriteString(w, "\n")
	}
}
