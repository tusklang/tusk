package ast

import "github.com/tusklang/tusk/tokenizer"

type AccessSpecifier interface {
	SetDecl(*ASTNode)
}

func parseAccessSpec(spec AccessSpecifier, lex []tokenizer.Token, i *int) error {
	*i++

	//match everything up to the next semicolon (that isn't enclosed in a brace)
	bm := braceMatcher(lex, i, []string{"{", "("}, []string{"}", ")"}, false, "terminator")
	g := grouper(bm)
	d, e := groupsToAST(g)

	if e != nil {
		return e
	}

	spec.SetDecl(d[0])
	return nil
}
