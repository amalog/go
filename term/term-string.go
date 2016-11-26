package term // import "github.com/amalog/go/term"

import "io"

type String string

func NewString(s string) String {
	return String(s)
}

func (s String) Format(w io.Writer, style Style) {
	io.WriteString(w, "\"")
	io.WriteString(w, string(s))
	io.WriteString(w, "\"")

	style.Terminate(w)
}

func (s String) Name() string {
	return `""`
}
