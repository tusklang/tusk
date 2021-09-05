package ast

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/tokenizer"
)

type IfStatement struct {
	Condition []*ASTNode
	Body      []*ASTNode
}

func (is *IfStatement) Parse(lex []tokenizer.Token, i *int) error {
	return ifwhileParse(is, lex, i)
}

func (is *IfStatement) SetCond(g []*ASTNode) {
	is.Condition = g
}

func (is *IfStatement) SetBody(g []*ASTNode) {
	is.Body = g
}

func (is *IfStatement) Type() string {
	return "if"
}

func (is *IfStatement) Compile(compiler *Compiler, class *types.StructType, node *ASTNode) constant.Constant {
	return nil
}
