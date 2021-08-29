package ast

import (
	"errors"

	"github.com/tusklang/tusk/tokenizer"
)

type Function struct {
	Name    string     //function name
	Params  []*ASTNode //parameter list
	RetType []*ASTNode //return type
	Body    []*ASTNode //function body
}

func (fh *Function) Parse(lex []tokenizer.Token, i *int) (e error) {

	if lex[*i].Type != "fn" {
		return errors.New("was not given a function")
	}

	*i++

	//read the return type
	//fn int main() {}
	//will also work, because if no braces are present, the next token is returned, and the brace matcher exits
	//if the next value is a variable name, then we know it's a void return type
	//so we will skip the return type

	if lex[*i].Type != "varname" {
		fh.RetType, e = groupsToAST(grouper(braceMatcher(lex, i, "(", ")", false, "")))
		if e != nil {
			return e
		}
	}

	*i++

	if lex[*i].Type != "varname" {
		return errors.New("function name was not provided, use the name '_' for an anonymous function")
	}

	fh.Name = lex[*i].Name

	*i++

	if lex[*i].Type != "(" { //it has to be a parenthesis for the paramlist
		return errors.New("functions require a parameter list")
	}

	fh.Params, e = groupsToAST(grouper(braceMatcher(lex, i, "(", ")", false, "")))

	if e != nil {
		return e
	}

	*i++

	fh.Body, e = groupsToAST(grouper(braceMatcher(lex, i, "{", "}", false, "terminator")))

	if e != nil {
		return e
	}

	return nil
}
