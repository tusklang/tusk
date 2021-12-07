package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type Pure struct {
	Declaration *ASTNode

	tok tokenizer.Token
}

func (p *Pure) Parse(lex []tokenizer.Token, i *int, stopAt []string) (e *errhandle.TuskError) {
	p.tok = lex[*i]
	return parseAccessSpec(p, lex, i)
}

func (p *Pure) SetDecl(node *ASTNode) {
	p.Declaration = node
}

func (p *Pure) GetMTok() tokenizer.Token {
	return p.tok
}

//cannot be compiled
func (p *Pure) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}
