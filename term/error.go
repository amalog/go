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

type ErrUnexpectedEof struct {
	Position scanner.Position
}

func (err *ErrUnexpectedEof) Error() string {
	return fmt.Sprintf("%s unexpected end of file", err.Position)
}

type Err struct {
	Token   *scanner.Token
	Message string
}

func (err *Err) Error() string {
	return fmt.Sprintf("%s %s", err.Token.Position, err.Message)
}
