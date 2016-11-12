package amalog // import "github.com/amalog/go"
import (
	"fmt"
	"io"
)

type Amalog struct {
	Out io.Writer
	Err io.Writer
}

func (ama *Amalog) Run(args []string) int {
	fmt.Fprintf(ama.Out, "TODO implement Amalog\n")
	return 1
}
