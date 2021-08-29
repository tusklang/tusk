package ast

import "github.com/tusklang/tusk/tokenizer"

func GenerateAST(tokens []tokenizer.Token) ([]*ASTNode, error) {
	g := grouper(tokens)
	a, e := groupsToAST(g)
	return a, e
}
