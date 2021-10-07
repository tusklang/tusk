package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Group interface {
	Parse([]tokenizer.Token, *int) error
	Compile(*Compiler, *data.Class, *ASTNode, *data.Function) data.Value
}
