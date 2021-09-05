package ast

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/tokenizer"
)

type DataType struct {
	Type tokenizer.Token
}

func (dt *DataType) Parse(lex []tokenizer.Token, i *int) error {
	dt.Type = lex[*i]
	return nil
}

func (dt *DataType) Compile(compiler *Compiler, class *types.StructType, node *ASTNode) constant.Constant {
	return nil
}
