package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Public struct {
	Declaration *ASTNode
}

func (p *Public) Parse(lex []tokenizer.Token, i *int, stopAt []string) (e error) {
	return parseAccessSpec(p, lex, i)
}

func (p *Public) SetDecl(node *ASTNode) {
	p.Declaration = node
}

//cannot be compiled
func (p *Public) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}
