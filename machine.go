package amalog // import "github.com/amalog/go"

import (
	"io"
	"os"

	"github.com/amalog/go/term"
)

type Machine struct {
	root term.Term // term which started it all
}

func NewMachine() *Machine {
	return &Machine{}
}

func (m *Machine) LoadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	return m.Load(file)
}

func (m *Machine) Load(r io.Reader) error {
	var err error
	m.root, err = term.ReadAll(r)
	return err
}

func (m *Machine) Call(name string, args ...term.Term) *Result {
	return nil
}
