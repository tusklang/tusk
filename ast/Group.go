package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/tokenizer"
)

type Group interface {
	Parse([]tokenizer.Token, *int) error
	Compile(*Compiler, *types.StructType, *ASTNode, *ir.Block) value.Value
}
