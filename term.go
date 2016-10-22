package prolog

import "math/big"

type Term interface {
}

type Variable interface {
}

type Number interface {
	AsBigRat() *big.Rat
}

func NewAtom(s string) Term {
	return nil
}

func IsGround(t Term) bool {
	return false
}

func TermFromBigRat(x *big.Rat) Term {
	return nil
}
