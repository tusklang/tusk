package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type DataType struct {
	Type tokenizer.Token
}

func (dt *DataType) Parse(lex []tokenizer.Token, i *int) error {
	dt.Type = lex[*i]
	return nil
}

func (dt *DataType) Compile(compiler *Compiler, class *types.StructType, node *ASTNode, block *ir.Block) data.Value {
	return nil
}
