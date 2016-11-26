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

// Arg returns the nth argument (1-based index) for a struct term. Panics
// if the term is not a struct or if the struct has too few arguments.
func Arg(t Term, n int) Term {
	if s, ok := t.(*Struct); ok {
		n-- // incoming index is 1-based
		if n < len(s.Args) {
			return s.Args[n]
		}
		panic("Term has too few arguments")
	}
	panic("Term is not a struct")
}

// Arity returns the length of the Args sequence for a struct term.  Otherwise,
// returns 0.
func Arity(t Term) int {
	if s, ok := t.(*Struct); ok {
		return len(s.Args)
	}
	return 0
}

// Clauses returns a slice of terms representing the clauses in the Db of a
// struct term.  For terms other than Struct, returns an empty slice.
func Clauses(t Term) []Term {
	if s, ok := t.(*Struct); ok {
		return []Term(s.Data)
	}
	return []Term{}
}
