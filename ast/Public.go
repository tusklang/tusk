package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Public struct {
	Declaration *ASTNode
}

func (p *Public) Parse(lex []tokenizer.Token, i *int) (e error) {
	return parseAccessSpec(p, lex, i)
}

func (p *Public) SetDecl(node *ASTNode) {
	p.Declaration = node
}

//cannot be compiled
func (p *Public) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {
	return nil
}
