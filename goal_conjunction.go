type GoalConjunction {
	Goals []Goal
}

func conjunction(goals ...Goal) Goal {
	return &GoalConjunction{goals}
}

func (self *GoalConjunction) Next(c Context) (bool,bool) {
	// TODO
	return false, false
}
