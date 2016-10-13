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

var one = big.NewRat(1,1)

// succ/2 is a built-in predicate
type Succ2 struct {
	X, Y term.Term	
}

func succ2(X, Y term.Term) Goal {
	return &Succ2{X,Y}
}

func (self *Succ2) Next() (bool, bool) {
	gx := term.Ground(self.X)
	gy := term.Ground(self.Y)
	if gx && gy {
		x := self.X.AsBigRat()
		y := self.Y.AsBigRat()
		z := new(big.Rat).Sub(y,x)
		return z.Cmp(one) == 0, false
	} else if gx {
		x := self.X.AsBigRat()
		y := new(big.Rat).Add(x,one)
		return Unify(self.Y,term.FromBigRat(y))
	} else if gy {
		y := self.Y.AsBigRat()
		x := new(big.Rat).Sub(y,one)
		return Unify(self.X,term.FromBigRat(x))
	} else {
		panic("mode error: X or Y should have been ground")
	}
}
