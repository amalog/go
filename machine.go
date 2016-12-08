package amalog // import "github.com/amalog/go"

import "github.com/amalog/go/term"

type Machine struct {
	root term.Term // term which started it all
}

func NewMachine() *Machine {
	return &Machine{}
}

func (m *Machine) Call(name string, args ...term.Term) Goal {
	n, err := term.NewAtom(name)
	if err != nil {
		panic(err)
	}
	s := term.NewSeq(args)
	t := &term.Struct{
		Name: n,
		Args: s,
	}
	return m.CallTerm(t)
}

func (m *Machine) CallTerm(goal term.Term) Goal {
	return &disjunction{goal}
}

func (m *Machine) Once(name string, args ...term.Term) Goal {
	goal := m.Call(name, args...)
	return &once{goal}
}
