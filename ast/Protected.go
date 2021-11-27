package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Protected struct {
	Declaration *ASTNode

	tok tokenizer.Token
}

func (p *Protected) Parse(lex []tokenizer.Token, i *int, stopAt []string) (e error) {
	p.tok = lex[*i]
	return parseAccessSpec(p, lex, i)
}

func (p *Protected) SetDecl(node *ASTNode) {
	p.Declaration = node
}

func (p *Protected) GetMTok() tokenizer.Token {
	return p.tok
}

//cannot be compiled
func (p *Protected) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}
