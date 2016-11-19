package term // import "github.com/amalog/go/term"

import (
	"io"
	"math/big"
	"reflect"
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

// Name returns a name for an arbitrary term.  It's easiest to understand
// by example:
//
//   * foo -> foo
//   * bar(x,y,z) -> bar
//   * # comment -> #
//   * "hi" -> ""
//
// A term implemented outside this package can choose its own name by
// implementing a `Name() string` method or having a public `Name string` field
// (if it's a struct). Otherwise, the name is "<unknown>".
func Name(t Term) string {
	type hasName interface {
		Name() string
	}
	if hasName, ok := t.(hasName); ok {
		return hasName.Name()
	}

	val := reflect.ValueOf(t)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Struct {
		if field := val.FieldByName("Name"); field.IsValid() {
			if field.Kind() == reflect.String {
				return field.String()
			}
		}
	}
	panic("bah!")
	return "<unknown>" // make compiler happy
}
