package ast

import (
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

func GenerateAST(tokens []tokenizer.Token) ([]*ASTNode, *errhandle.TuskError) {
	g, e := grouper(tokens)

	if e != nil {
		return nil, e
	}

	a, e := groupsToAST(g)
	return a, e
}
