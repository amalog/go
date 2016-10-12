/*
%% succ(num,num).
%% succ(+X,+Y) is semidet.
%% succ(+X,-Y) is det.
%% succ(-X,+Y) is det.
succ(X,Y) :-
	Y #= X + 1.

In the function names below, "i" represents an input argument
and "o" represents an output argument.
*/

import (
	"math/big"
)

var one = big.NewRat(1,1)

func succ2_ii_semidet(X, Y term.Term) bool {
	x := X.AsBigRat()
	y := Y.AsBigRat()
	z := new(big.Rat).Sub(y,x)
	return z.Cmp(one) == 0
}

func succ2_io_det(X term.Term, Y term.Variable) {
	x := X.AsBigRat()
	y := new(big.Rat).Add(x,one)
	if !Y.Unify(term.FromBigRat(y)) {
		panic("mode error: Y should have been an unbound variable")
	}
}

func succ2_oi_det(X term.Variable, Y term.Term) {
	y := Y.AsBigRat()
	x := new(big.Rat).Sub(y,one)
	if !X.Unify(term.FromBigRat(x)) {
		panic("mode error: X should have been an unbound variable")
	}
}