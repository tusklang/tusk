package ast

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/tokenizer"
)

type Protected struct {
	Declaration *ASTNode
}

func (p *Protected) Parse(lex []tokenizer.Token, i *int) (e error) {
	return parseAccessSpec(p, lex, i)
}

func (p *Protected) SetDecl(node *ASTNode) {
	p.Declaration = node
}

//cannot be compiled
func (p *Protected) Compile(compiler *Compiler, class *types.StructType, node *ASTNode) constant.Constant {
	return nil
}
