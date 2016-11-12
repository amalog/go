package term // import "github.com/amalog/go/term"

import (
	"io"
	"math/big"
)

type Term interface {
	// Format writes a textual representation of this term according to the given
	// style.
	Format(io.Writer, Style)

	// String produces the canonical representation of this term when it
	// stands alone.
	String() string
}

type Number *big.Rat

func IsGround(t Term) bool {
	return false
}

func TermFromBigRat(x *big.Rat) Term {
	return nil
}
