package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type Group interface {
	Parse([]tokenizer.Token, *int, []string) *errhandle.TuskError
	GetMTok() tokenizer.Token
	Compile(*Compiler, *data.Class, *ASTNode, *data.Function) data.Value
}
