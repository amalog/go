package scanner // import "github.com/amalog/go/scanner"
import "fmt"

type Position struct {
	Filename string
	Line     int
	Column   int
}

func (pos Position) String() string {
	f := pos.Filename
	if f == "" {
		f = "<input>"
	}
	return fmt.Sprintf("%s:%d:%d", f, pos.Line, pos.Column)
}
