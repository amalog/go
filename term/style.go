package term // import "github.com/amalog/go/term"
import "io"

type Style struct {
	// OmitNewline should be to set to true to avoid including newlines after each
	// term.
	OmitNewline bool

	// OmitTerminator should be set to true to avoid including term terminators in
	// the output.
	OmitTerminator bool

	// Indent indicates the number of levels of indentation that should be applied
	// before the term is displayed.
	Indent int
}

// WriteIndent writes the appropriate amount of indentation for this style.
func (s Style) WriteIndent(w io.Writer) {
	for i := 0; i < s.Indent; i++ {
		io.WriteString(w, "    ")
	}
}

// Terminate writes the appropriate term terminator for this style.
func (s Style) Terminate(w io.Writer) {
	if s.OmitNewline {
		if !s.OmitTerminator {
			io.WriteString(w, ",")
		}
	} else {
		io.WriteString(w, "\n")
	}
}
