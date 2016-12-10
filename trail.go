package amalog // import "github.com/amalog/go"

import "github.com/amalog/go/term"

type Barrier int

type Trail interface {
	BacktrackTo(Barrier)
	DropBarrier(Barrier)
	PushBarrier() Barrier
}

func unify(X, Y term.Term) Goal {
	return nil
}

func Unify(c Context, x, y term.Term) bool {
	return true
}
