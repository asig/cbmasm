package errors

import (
	"fmt"
	"github.com/asig/cbmasm/pkg/text"
)

type Error struct {
	Pos text.Pos
	Msg string
}

func (e Error) String() string {
	return fmt.Sprintf("%s, line %d, col %d: %s", e.Pos.Filename, e.Pos.Line, e.Pos.Col, e.Msg)
}

type Sink interface {
	AddError(pos text.Pos, message string);
}

