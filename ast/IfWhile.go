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

	cg := grouper(braceMatcher(lex, i, "(", ")", false, ""))
	ca, e := groupsToAST(cg)
	if e != nil {
		return e
	}
	statement.SetCond(ca)

	*i++

	bg := grouper(braceMatcher(lex, i, "{", "}", false, "terminator"))
	ba, e := groupsToAST(bg)
	if e != nil {
		return e
	}
	statement.SetBody(ba)

	return nil
}

type IfStatement struct {
	Condition []*ASTNode
	Body      []*ASTNode
}

func (is *IfStatement) Parse(lex []tokenizer.Token, i *int) error {
	return ifwhileParse(is, lex, i)
}

func (is *IfStatement) SetCond(g []*ASTNode) {
	is.Condition = g
}

func (is *IfStatement) SetBody(g []*ASTNode) {
	is.Body = g
}

func (is *IfStatement) Type() string {
	return "if"
}

type WhileStatement struct {
	Condition []*ASTNode
	Body      []*ASTNode
}

func (ws *WhileStatement) Parse(lex []tokenizer.Token, i *int) error {
	return ifwhileParse(ws, lex, i)
}

func (ws *WhileStatement) SetCond(g []*ASTNode) {
	ws.Condition = g
}

func (ws *WhileStatement) SetBody(g []*ASTNode) {
	ws.Body = g
}

func (ws *WhileStatement) Type() string {
	return "while"
}
