package amalog // import "github.com/amalog/go"

import (
	"log"

	"github.com/amalog/go/term"
)

// conjunction searches for solutions which satisfy all goals within a body
type conjunction struct {
	body term.Db
}

func (self *conjunction) Next(c Context) (bool, bool) {
	return true, false
}
