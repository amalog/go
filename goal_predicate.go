type PredicateGoal struct {
	Predicate PredicateId
	Args      []term.Term

	clauseIndex int
	clauses     []Goal
}

func (self *PredicateGoal) Next(c Context) (bool, bool) {
	if self.clauses == nil {
		self.clauses = lookupClauses(self.Predicate, self.Args)
	}

	// look for solutions from each clause in turn
	for {
		if clauseIndex >= len(self.clauses) {
			// already investigated all clauses
			return false, false
		}
		clause := self.clauses[self.clauseIndex]
		ok, more := clause.Next(c)
		if ok {
			if more {
				return true, true
			} else {
				self.clauseIndex++
				return true, self.clauseIndex < len(self.clauses)
			}
		} else {
			// look for solutions in the next clause
			self.clauseIndex++
		}
	}
}

func (*PredicateGoal) Cleanup() {}
