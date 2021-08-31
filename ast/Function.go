package ast

import (
	"errors"

	"github.com/tusklang/tusk/tokenizer"
)

type Parameter struct {
	Name string
	Type *ASTNode
}

type Function struct {
	Name    string      //function name
	Params  []Parameter //parameter list
	RetType *ASTNode    //return type
	Body    []*ASTNode  //function body
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
		rt, e := groupsToAST(groupSpecific(lex, 1, i))
		fh.RetType = rt[0]
		if e != nil {
			return e
		}
	}

	if lex[*i].Type == "varname" {
		fh.Name = lex[*i].Name
		*i++
	}

	if lex[*i].Type != "(" { //it has to be a parenthesis for the paramlist
		return errors.New("functions require a parameter list")
	}

	p, e := groupsToAST(grouper(braceMatcher(lex, i, []string{"("}, []string{")"}, false, "")))
	sub := p[0].Group.(*Block).Sub
	plist := make([]Parameter, len(sub))

	for k, v := range sub {

		switch g := v.Group.(type) {
		case *Operation:
			if g.OpType != ":" {
				return errors.New("invalid syntax: named parameters must have a type")
			}

			plist[k] = Parameter{
				Name: v.Left[0].Group.(*DataValue).Value.Name,
				Type: v.Right[0],
			}

		case *DataType:
			plist[k] = Parameter{
				Type: v,
			}

		//a classname
		case *DataValue:
			plist[k] = Parameter{
				Type: v,
			}

		}
		_ = v
	}

	fh.Params = plist

	if e != nil {
		return e
	}

	*i++

	fh.Body, e = groupsToAST(grouper(braceMatcher(lex, i, []string{"{"}, []string{"}"}, false, "terminator")))

	if e != nil {
		return e
	}

	return nil
}
