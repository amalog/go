package term // import "github.com/amalog/go/term"

import "fmt"

type String string

func NewString(s string) String {
	return String(s)
}

func (a String) String() string {
	return fmt.Sprintf("\"%s\";\n", string(a))
}
