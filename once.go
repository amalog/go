package amalog // import "github.com/amalog/go"

// once wraps a goal to give a new goal which only produces the first solution.
type once struct {
	g Goal
}

func (self *once) Next(c Context) (bool, bool) {
	ok := false
	if self.g != nil {
		ok, _ = self.g.Next(c)
		self.g = nil // after first solution, we no longer need the goal
	}

	return ok, false
}
