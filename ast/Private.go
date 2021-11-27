package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type Private struct {
	Declaration *ASTNode

	tok tokenizer.Token
}

func (p *Private) Parse(lex []tokenizer.Token, i *int, stopAt []string) (e *errhandle.TuskError) {
	p.tok = lex[*i]
	return parseAccessSpec(p, lex, i)
}

func (p *Private) SetDecl(node *ASTNode) {
	p.Declaration = node
}

func (p *Private) GetMTok() tokenizer.Token {
	return p.tok
}

//cannot be compiled
func (p *Private) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}
