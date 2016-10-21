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
	X, Y term.Term
}

func succ2(X, Y term.Term) Goal {
	return &Succ2{X, Y}
}

func (self *Succ2) Next(c Context) (bool, bool) {
	if term.IsGround(self.X) {
		x := self.X.AsBigRat()
		y := new(big.Rat).Add(x, one)
		return Unify(c, self.Y, term.FromBigRat(y))
	} else if term.IsGround(self.Y) {
		y := self.Y.AsBigRat()
		x := new(big.Rat).Sub(y, one)
		return Unify(c, self.X, term.FromBigRat(x))
	} else {
		panic("mode error: X or Y should have been ground")
	}
}

func (*Succ2) Cleanup() {}
