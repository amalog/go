package amalog // import "github.com/amalog/go"

import "github.com/amalog/go/term"

// disjunction searches for all solutions to a given goal term
type disjunction struct {
	goal term.Term
}

func (self *disjunction) Next(c Context) (bool, bool) {
	return false, false
}
