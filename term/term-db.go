package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"io"
)

type Db []Term

func NewDb(args []Term) Db {
	return Db(args)
}

func (s Db) String() string {
	buf := new(bytes.Buffer)
	s.Format(buf, Style{})
	return buf.String()
}

func (s Db) Format(w io.Writer, style Style) {
	for _, t := range s {
		style.WriteIndent(w)
		t.Format(w, style)
	}
}
