package ast

import (
	"errors"

	"github.com/tusklang/tusk/tokenizer"
)

type IfWhile interface {
	SetCond([]*ASTNode)
	SetBody([]*ASTNode)
	Type() string
	SetSTok(tokenizer.Token)
	SetCondTok(tokenizer.Token)
	SetBTok(tokenizer.Token)
}

func ifwhileParse(statement IfWhile, lex []tokenizer.Token, i *int) error {

	if lex[*i].Type != statement.Type() {
		return errors.New("token given does not match parse expectation")
	}

	statement.SetSTok(lex[*i])

	*i++

	statement.SetCondTok(lex[*i])
	cg := grouper(braceMatcher(lex, i, []string{"("}, []string{")"}, false, ""))
	ca, e := groupsToAST(cg)
	if e != nil {
		return e
	}
	statement.SetCond(ca)

	*i++

	statement.SetBTok(lex[*i])
	bg := grouper(braceMatcher(lex, i, []string{"{"}, []string{"}"}, false, "terminator"))
	ba, e := groupsToAST(bg)
	if e != nil {
		return e
	}
	statement.SetBody(ba)

	switch statement.(type) {
	case *IfStatement:
		//if it's an if statement, we need to check if there is an `else` clause
		if *i+1 < len(lex) && lex[*i+1].Name == "else" {
			//else clause detected
			*i += 2 //skip the semicolon & "else"
			elsebody := grouper(braceMatcher(lex, i, []string{"{"}, []string{"}"}, false, "terminator"))
			statement.(*IfStatement).ElseBody, e = groupsToAST(elsebody)

			if e != nil {
				return e
			}
		}
	}

	return nil
}
