package ast

import (
	"github.com/tusklang/tusk/errhandle"
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

func ifwhileParse(statement IfWhile, lex []tokenizer.Token, i *int) *errhandle.TuskError {

	statement.SetSTok(lex[*i])

	*i++

	statement.SetCondTok(lex[*i])
	cgbm, e := braceMatcher(lex, i, []string{"("}, []string{")"}, false, "")
	if e != nil {
		return e
	}
	cg, e := grouper(cgbm)
	if e != nil {
		return e
	}
	ca, e := groupsToAST(cg)
	if e != nil {
		return e
	}
	statement.SetCond(ca)

	*i++

	statement.SetBTok(lex[*i])
	bgbm, e := braceMatcher(lex, i, []string{"{"}, []string{"}"}, false, "terminator")
	if e != nil {
		return e
	}
	bg, e := grouper(bgbm)
	if e != nil {
		return e
	}
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
			elsebodybm, e := braceMatcher(lex, i, []string{"{"}, []string{"}"}, false, "terminator")
			if e != nil {
				return e
			}
			elsebody, e := grouper(elsebodybm)
			if e != nil {
				return e
			}
			statement.(*IfStatement).ElseBody, e = groupsToAST(elsebody)

			if e != nil {
				return e
			}
		}
	}

	return nil
}
