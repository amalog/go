package term // import "github.com/amalog/go/term"
import "bytes"

type Seq []Term

func NewSeq(args []Term) Seq {
	return Seq(args)
}

func (s Seq) String() string {
	buf := new(bytes.Buffer)
	final := len(s) - 1
	for i, t := range s {
		buf.WriteString(t.String())
		if i < final {
			buf.WriteString(", ")
		}
	}

	return buf.String()
}
