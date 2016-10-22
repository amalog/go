package prolog

/*
%% succ(num,num).
%% succ(+,+) is semidet.
%% succ(+,-) is det.
%% succ(-,+) is det.
succ(X,Y) :-
	Y #= X + 1.
*/

import (
	"math/big"
)

var one = big.NewRat(1, 1)

// succ/2 is a built-in predicate
type Succ2 struct {
	X, Y Term
}

func succ2(X, Y Term) Goal {
	return &Succ2{X, Y}
}

func (self *Succ2) Next(c Context) (bool, bool) {
	if IsGround(self.X) {
		if x, ok := self.X.(Number); ok {
			y := new(big.Rat).Add(x.AsBigRat(), one)
			return Unify(c, self.Y, TermFromBigRat(y))
		} else {
			panic("type error: X is not a number")
		}
	} else if IsGround(self.Y) {
		if y, ok := self.Y.(Number); ok {
			x := new(big.Rat).Sub(y.AsBigRat(), one)
			return Unify(c, self.X, TermFromBigRat(x))
		} else {
			panic("type error: Y is not a number")
		}
	} else {
		panic("mode error: X or Y should have been ground")
	}
}
