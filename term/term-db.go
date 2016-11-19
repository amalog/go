package term // import "github.com/amalog/go/term"

import "io"

type Db []Term

func NewDb(args []Term) Db {
	return Db(args)
}

func (s Db) Format(w io.Writer, style Style) {
	for _, t := range s {
		style.WriteIndent(w)
		t.Format(w, style)
	}
}
