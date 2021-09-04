package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/tokenizer"
)

type DataValue struct {
	Value tokenizer.Token
}

func (dv *DataValue) Parse(lex []tokenizer.Token, i *int) error {
	dv.Value = lex[*i]
	return nil
}

func (dv *DataValue) Compile(compiler *Compiler, class *types.StructType, node *ASTNode) {}
