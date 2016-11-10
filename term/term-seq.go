package term // import "github.com/amalog/go/term"

import "fmt"

type Seq []Term

func NewSeq(args []Term) Seq {
	return Seq(args)
}

func (s Seq) String() string {
	return fmt.Sprintf("%#v;\n", s)
}
