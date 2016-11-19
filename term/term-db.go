package term // import "github.com/amalog/go/term"

import "io"

type Db []Term

func NewDb(args []Term) Db {
	return Db(args)
}

func (s Db) Format(w io.Writer, style Style) {
	isRoot := style.IsRoot
	style.IsRoot = false // our children are not root

	prevName := "#" // won't impose extra newlines
	for _, t := range s {
		if isRoot {
			name := Name(t)
			if prevName != "#" && prevName != name {
				io.WriteString(w, "\n")
			}
			prevName = name
		}

		style.WriteIndent(w)
		t.Format(w, style)
	}
}
