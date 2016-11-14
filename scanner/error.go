package scanner // import "github.com/amalog/go/scanner"
import "fmt"

type SyntaxError struct {
	Position Position
	Message  string
}

func (err *SyntaxError) Error() string {
	return fmt.Sprintf("%s: %s", err.Position, err.Message)
}
