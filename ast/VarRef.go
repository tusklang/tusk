package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/tokenizer"
)

type VarRef struct {
	Name string
}

func (vr *VarRef) Parse(lex []tokenizer.Token, i *int) error {
	vr.Name = lex[*i].Name
	return nil
}

func (vr *VarRef) Compile(compiler *Compiler, class *types.StructType, node *ASTNode, block *ir.Block) value.Value {
	return nil
}
