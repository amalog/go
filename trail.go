package prolog

type Barrier int

type Trail interface {
	BacktrackTo(Barrier)
	DropBarrier(Barrier)
	PushBarrier() Barrier
}

func unify(X, Y Term) Goal {
	return nil
}

func Unify(c Context, x, y Term) (bool, bool) {
	return false, false
}
