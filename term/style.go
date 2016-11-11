package term // import "github.com/amalog/go/term"

type Style struct {
	// OmitNewline should be to set to true to avoid including newlines after each
	// term.
	OmitNewline bool

	// OmitTerminator should be set to true to avoid including term terminators in
	// the output.
	OmitTerminator bool
}
