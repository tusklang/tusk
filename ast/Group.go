package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/tokenizer"
)

type Group interface {
	Parse([]tokenizer.Token, *int) error
	Compile(*types.StructType, *ASTNode)
}
