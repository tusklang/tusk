package ast

import (
	"github.com/tusklang/tusk/tokenizer"
)

type Public struct {
	Declaration *ASTNode
}

func (p *Public) Parse(lex []tokenizer.Token, i *int) (e error) {
	*i++

	//match everything up to the next semicolon (that isn't enclosed in a brace)
	bm := braceMatcher(lex, i, []string{"{", "("}, []string{"}", ")"}, false, "terminator")
	g := grouper(bm)
	d, e := groupsToAST(g)
	p.Declaration = d[0]
	return
}
