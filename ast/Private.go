package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Private struct {
	Declaration *ASTNode
}

func (p *Private) Parse(lex []tokenizer.Token, i *int) (e error) {
	return parseAccessSpec(p, lex, i)
}

func (p *Private) SetDecl(node *ASTNode) {
	p.Declaration = node
}

//cannot be compiled
func (p *Private) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}
