package term // import "github.com/amalog/go/term"

import (
	"fmt"
	"io"
	"regexp"
)

type Var struct {
	Name  string
	Value Term // nil to indicate unbound variable
}

var consecutiveCapitals *regexp.Regexp

func init() {
	consecutiveCapitals = regexp.MustCompile(`[A-Z][A-Z]`)
}

func NewVar(name string) (*Var, error) {
	if consecutiveCapitals.MatchString(name) {
		err := fmt.Errorf("variable name (%s) may not have consecutive uppercase letters", name)
		return nil, err
	}
	return &Var{Name: name}, nil
}

func (v *Var) Format(w io.Writer, style Style) {
	io.WriteString(w, v.Name)

	style.Terminate(w)
}
