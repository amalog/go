package amalog // import "github.com/amalog/go"

import (
	"os"

	"github.com/amalog/go/term"
)

type Machine struct {
	root term.Term // term which started it all
}

func NewMachine() *Machine {
	return &Machine{}
}

// LoadRoot populates the machine's root term with the contents of a file
// and loads all of that term's module dependencies.
func (m *Machine) LoadRoot(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	m.root, err = term.ReadAll(file)
	if err != nil {
		return err
	}

	return nil
}

func (m *Machine) Call(name string, args ...term.Term) *Result {
	return nil
}
