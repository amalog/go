package prolog

type GoalConjunction struct {
	Goals []Goal

	choicepoints []choicepoint
}

type choicepoint struct {
	index   int
	barrier Barrier
}

func conjunction(goals ...Goal) Goal {
	return &GoalConjunction{Goals: goals}
}

func (self *GoalConjunction) Next(c Context) (bool, bool) {
	goalIndex := 0
	if len(self.choicepoints) > 0 {
		p := self.choicepoints[len(self.choicepoints)-1]
		self.choicepoints = self.choicepoints[0 : len(self.choicepoints)-1]
		goalIndex = p.index
		c.Trail().BacktrackTo(p.barrier)
		c.Trail().DropBarrier(p.barrier)
	}

	for goalIndex < len(self.Goals) {
		goal := self.Goals[goalIndex]
		barrier := c.Trail().PushBarrier()
		ok, more := goal.Next(c)
		if ok {
			if more {
				// record a choicepoint
				p := choicepoint{goalIndex, barrier}
				self.choicepoints = append(self.choicepoints, p)
			}
			goalIndex++
			continue
		} else {
			// TODO backtrack to previous choicepoint
		}
	}

	more := len(self.choicepoints) > 0
	if !more {
		for _, goal := range self.Goals {
			if x, ok := goal.(NeedsCleanup); ok {
				x.Cleanup()
			}
		}
	}
	return true, more
}
