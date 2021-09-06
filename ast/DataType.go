package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/tokenizer"
)

type DataType struct {
	Type tokenizer.Token
}

func (dt *DataType) Parse(lex []tokenizer.Token, i *int) error {
	dt.Type = lex[*i]
	return nil
}

func (dt *DataType) Compile(compiler *Compiler, class *types.StructType, node *ASTNode, block *ir.Block) value.Value {
	switch dt.Type.Name {

	//return a value with the type of the data
	case "i32":
		return constant.NewInt(types.I32, 0)
	default:
		return nil
	}
}
