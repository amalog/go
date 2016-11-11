package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"io"
)

type Seq []Term

func NewSeq(args []Term) Seq {
	return Seq(args)
}

func (s Seq) String() string {
	buf := new(bytes.Buffer)
	s.Format(buf, Style{})
	return buf.String()
}

func (s Seq) Format(w io.Writer, style Style) {
	style.OmitTerminator = true
	style.OmitNewline = true

	final := len(s) - 1
	for i, t := range s {
		t.Format(w, style)
		if i < final {
			io.WriteString(w, ", ")
		}
	}
}
