package amalog // import "github.com/amalog/go"
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/amalog/go/term"
)

type Amalog struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

// TODO most of this should be rewritten in Amalog so that it doesn't
// TODO have to be recreated in each implementation.
func (ama *Amalog) Run(args []string) int {
	if len(args) == 0 {
		return ama.CmdRepl()
	}

	// run an Amalog script
	if _, err := os.Stat(args[0]); err == nil {
		return ama.CmdRun(args[0])
	}

	// subcommands
	switch args[0] {
	case "format":
		if len(args) < 2 {
			fmt.Fprintln(ama.Err, "format command needs a file argument")
			return 1
		}
		return ama.CmdFormat(args[1])
	case "run":
		if len(args) < 2 {
			fmt.Fprintln(ama.Err, "run command needs a file argument")
			return 1
		}
		return ama.CmdRun(args[1])
	default:
		fmt.Fprintf(ama.Err, "Unrecognized command: %s", args[0])
		return 1
	}
}

func (ama *Amalog) CmdRepl() int {
	buf := bufio.NewReader(ama.In)

	style := term.Style{}
	for {
		fmt.Fprintf(ama.Out, "?- ")
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(ama.Err, err)
		}

		reader := term.NewReader(strings.NewReader(line))
		for {
			t, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintln(ama.Err, err)
				break
			}
			t.Format(ama.Out, style)
		}
	}
	return 0
}

func (ama *Amalog) CmdRun(filename string) int {
	m := NewMachine()
	err := m.LoadRoot(filename)
	if err != nil {
		fmt.Fprintf(ama.Err, "%s\n", err)
		return 1
	}
	ok, _ := m.Once("main", World).Next(nil)
	if !ok {
		fmt.Fprintf(ama.Err, "main/0 failed")
		return 1
	}
	return 0
}

func (ama *Amalog) CmdFormat(filename string) int {
	// open Amalog source code file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(ama.Err, "%s\n", err)
		return 1
	}

	// read and output terms
	t, err := term.ReadAllAsTerm(file)
	if err != nil {
		fmt.Fprintf(ama.Err, "read: %s", err)
		return 1
	}
	style := term.Style{IsRoot: true}
	t.Format(ama.Out, style)
	return 0
}
