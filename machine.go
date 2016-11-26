package amalog // import "github.com/amalog/go"

import "github.com/amalog/go/term"

type Machine struct {
	root term.Term // term which started it all
}

func NewMachine() *Machine {
	return &Machine{}
}

func (m *Machine) Call(name string, args ...term.Term) *Result {
	return nil
}
