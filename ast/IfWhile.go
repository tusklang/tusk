package ast

import (
	"errors"

	"github.com/tusklang/tusk/tokenizer"
)

type IfWhile interface {
	SetCond([]*ASTNode)
	SetBody([]*ASTNode)
	Type() string
}

func ifwhileParse(statement IfWhile, lex []tokenizer.Token, i *int) error {

	if lex[*i].Type != statement.Type() {
		return errors.New("token given does not match parse expectation")
	}

	*i++

	cg := grouper(braceMatcher(lex, i, []string{"("}, []string{")"}, false, ""))
	ca, e := groupsToAST(cg)
	if e != nil {
		return e
	}
	statement.SetCond(ca)

	*i++

	bg := grouper(braceMatcher(lex, i, []string{"{"}, []string{"}"}, false, "terminator"))
	ba, e := groupsToAST(bg)
	if e != nil {
		return e
	}
	statement.SetBody(ba)

	return nil
}
