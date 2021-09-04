package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/tokenizer"
)

type VarRef struct {
	Name string
}

func (vr *VarRef) Parse(lex []tokenizer.Token, i *int) error {
	vr.Name = lex[*i].Name
	return nil
}

func (vr *VarRef) Compile(class *types.StructType, node *ASTNode) {}
