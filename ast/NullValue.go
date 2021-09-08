package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type NullValue struct{}

func (nv *NullValue) Parse(lex []tokenizer.Token, i *int) error {
	return nil
}

func (dv *NullValue) Compile(compiler *Compiler, class *types.StructType, node *ASTNode, block *ir.Block) data.Value {
	return nil
}
