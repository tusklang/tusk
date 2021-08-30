package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type Public struct {
	Declaration *ASTNode
}

func (p *Public) Parse(lex []tokenizer.Token, i *int) (e error) {
	*i++
	g := groupSpecific(lex, 1, i)
	d, e := groupsToAST(g)
	p.Declaration = d[0]
	*i-- //decrement because the outer grouper function already increments once
	return
}
