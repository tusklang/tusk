package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Public struct {
	Declaration *ASTNode

	tok tokenizer.Token
}

func (p *Public) Parse(lex []tokenizer.Token, i *int, stopAt []string) (e error) {
	p.tok = lex[*i]
	return parseAccessSpec(p, lex, i)
}

func (p *Public) SetDecl(node *ASTNode) {
	p.Declaration = node
}

func (p *Public) GetMTok() tokenizer.Token {
	return p.tok
}

//cannot be compiled
func (p *Public) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}
