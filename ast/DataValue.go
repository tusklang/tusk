package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/tokenizer"
)

type DataValue struct {
	Value value.Value
}

func (dv *DataValue) Parse(lex []tokenizer.Token, i *int) error {
	dv.Value = getValue(lex[*i]) //convert the lex into an llvm value
	return nil
}

func (dv *DataValue) Compile(compiler *Compiler, class *types.StructType, node *ASTNode, block *ir.Block) value.Value {
	return dv.Value
}
