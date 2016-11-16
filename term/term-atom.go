package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
)

type Atom string

var okAtomName *regexp.Regexp

func init() {
	okAtomName = regexp.MustCompile(`^[a-z][a-z_]*$`)
}

func NewAtom(s string) (Atom, error) {
	if okAtomName.MatchString(s) {
		return Atom(s), nil
	}

	err := fmt.Errorf("invalid atom name: %s", s)
	return "", err
}

func (a Atom) String() string {
	buf := new(bytes.Buffer)
	a.Format(buf, Style{})
	return buf.String()
}

func (a Atom) Format(w io.Writer, style Style) {
	io.WriteString(w, string(a))
	style.Terminate(w)
}

func (a Atom) Name() string {
	return string(a)
}
