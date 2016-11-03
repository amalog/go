package term // import "github.com/amalog/go/term"

import "math/big"

type Term interface {
}

type Variable interface {
}

type Number interface {
	AsBigRat() *big.Rat
}

type Atom interface {
}

type Seq interface {
}

type Database interface {
}

type Struct struct {
	Context Variable
	Name    Atom
	Args    Seq
	Data    Database
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
