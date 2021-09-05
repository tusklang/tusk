package ast

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/tokenizer"
)

type Operation struct {
	OpType string
	Token  *tokenizer.Token
}

func (o *Operation) Parse(lex []tokenizer.Token, i *int) error {

	o.Token = &lex[*i]
	o.OpType = lex[*i].Name

	return nil
}

func (o *Operation) Compile(compiler *Compiler, class *types.StructType, node *ASTNode) constant.Constant {
	return nil
}
