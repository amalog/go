package amalog // import "github.com/amalog/go"

import (
	"log"

	"github.com/amalog/go/term"
)

// disjunction searches for all solutions to a given goal term
type disjunction struct {
	db   term.Db   // database within which to find clauses
	goal term.Term // clauses should match this goal

	cursor      *term.Cursor // where to resume searching for clauses
	body        Goal         // body of the current matching clause, if any
	moreClauses bool         // are there more clauses beyond the current one?
}

func (self *disjunction) Next(c Context) (bool, bool) {
	// look for more solutions from the body
	if self.body != nil {
		ok, moreBodySolutions := self.body.Next(c)
		if !ok || !moreBodySolutions {
			self.body = nil // done with this body
		}
		if ok {
			return true, (moreBodySolutions || self.moreClauses)
		}
	}

	// search for matching clauses
	if self.cursor == nil {
		self.cursor = self.db.NewQuery(self.goal).Run()
	}
	for {
		var clause term.Term
		clause, self.moreClauses = self.cursor.Next()
		if clause == nil {
			break
		}

		if Unify(c, self.goal, clause) {
			body, err := term.Body(clause)
			if err != nil { // clause has no body. succeed
				return true, self.moreClauses
			}
			self.body = &conjunction{body}
			return self.Next(c)
		}

		if !self.moreClauses {
			break
		}
	}

	return false, false
}
