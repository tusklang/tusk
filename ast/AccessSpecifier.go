package ast

import (
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type AccessSpecifier interface {
	SetDecl(*ASTNode)
}

func parseAccessSpec(spec AccessSpecifier, lex []tokenizer.Token, i *int) *errhandle.TuskError {
	*i++

	//match everything up to the next semicolon (that isn't enclosed in a brace)
	bm := braceMatcher(lex, i, allopeners, allclosers, false, "terminator")
	g, e := grouper(bm)

	if e != nil {
		return e
	}

	d, e := groupsToAST(g)

	if e != nil {
		return e
	}

	spec.SetDecl(d[0])
	return nil
}
