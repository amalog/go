package term // import "github.com/amalog/go/term"
import "fmt"

type Var struct {
	Name  string
	Value *Term
}

func (v *Var) String() string {
	return fmt.Sprintf("%s,\n", v.Name)
}
