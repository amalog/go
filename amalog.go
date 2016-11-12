package amalog // import "github.com/amalog/go"
import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/amalog/go/term"
)

type Amalog struct {
	Out io.Writer
	Err io.Writer
}

func (ama *Amalog) Run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(ama.Err, "Usage: ama foo.ama\n")
		return 1
	}

	return ama.CmdFormat(args[0])
}

func (ama *Amalog) CmdFormat(filename string) int {
	// open Amalog source code file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(ama.Err, "%s\n", err)
		return 1
	}
	buf := bufio.NewReader(file)

	// read and output terms
	style := term.Style{}
	reader := term.NewReader(buf)
	for {
		t, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(ama.Err, "read: %s", err)
			return 1
		}

		t.Format(ama.Out, style)
	}

	return 0
}
