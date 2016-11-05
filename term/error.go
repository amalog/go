package term // import "github.com/amalog/go/term"

import (
	"fmt"

	"github.com/amalog/go/scanner"
)

type ErrUnexpectedToken struct {
	Token *scanner.Token
}

func (err *ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("%s unexpected token: %s", err.Token.Position, err.Token)
}
