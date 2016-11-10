package term // import "github.com/amalog/go/term"

import "math/big"

type Term interface {
	// String produces the canonical representation of this Amalog term.
	String() string
}

type Number *big.Rat

type Database []Term

func IsGround(t Term) bool {
	return false
}

func TermFromBigRat(x *big.Rat) Term {
	return nil
}
