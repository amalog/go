package term // import "github.com/amalog/go/term"

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

type Comment string

func NewComment(s string) (Comment, error) {
	if strings.ContainsRune(s, '\n') {
		return "", errors.New("comments may not contain a newline")
	}
	return Comment(s), nil
}

func (s Comment) String() string {
	buf := new(bytes.Buffer)
	s.Format(buf, Style{})
	return buf.String()
}

func (s Comment) Format(w io.Writer, style Style) {
	io.WriteString(w, "#")
	io.WriteString(w, string(s))
	io.WriteString(w, "\n")
}
