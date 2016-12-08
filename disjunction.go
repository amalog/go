package amalog // import "github.com/amalog/go"

import "github.com/amalog/go/term"

// disjunction searches for all solutions to a given goal term
type disjunction struct {
	root term.Term
	goal term.Term
}

func (self *disjunction) Next(c Context) (bool, bool) {
	/*
		cursor := self.root.Data.NewQuery(self.goal).Run()
		for {
			clause, more := cursor.Next()
			if clause == nil {
				break
			}

			// TODO unify goal with clause
			// TODO if success, construct conjunction of goal.Data
			// TODO iterate Next() on conjunction

			if !more {
				break
			}
		}
	*/
	return false, false
}
