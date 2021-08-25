package grouper

import (
	"errors"

	"github.com/tusklang/tusk/tokenizer"
)

type IfWhile interface {
	SetCond([]Group)
	SetBody([]Group)
	Type() string
}

func ifwhileParse(statement IfWhile, lex []tokenizer.Token, i *int) error {

	if lex[*i].Type != statement.Type() {
		return errors.New("token given does not match parse expectation")
	}

	*i++

	statement.SetCond(Grouper(braceMatcher(lex, i, "(", ")", false, "")))

	*i++

	statement.SetBody(Grouper(braceMatcher(lex, i, "{", "}", false, "terminator")))

	return nil
}

type IfStatement struct {
	Condition []Group
	Body      []Group
}

func (is *IfStatement) Parse(lex []tokenizer.Token, i *int) error {
	return ifwhileParse(is, lex, i)
}

func (is *IfStatement) SetCond(g []Group) {
	is.Condition = g
}

func (is *IfStatement) SetBody(g []Group) {
	is.Body = g
}

func (is *IfStatement) Type() string {
	return "if"
}

type WhileStatement struct {
	Condition []Group
	Body      []Group
}

func (ws *WhileStatement) Parse(lex []tokenizer.Token, i *int) error {
	return ifwhileParse(ws, lex, i)
}

func (ws *WhileStatement) SetCond(g []Group) {
	ws.Condition = g
}

func (ws *WhileStatement) SetBody(g []Group) {
	ws.Body = g
}

func (ws *WhileStatement) Type() string {
	return "while"
}