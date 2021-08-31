package ast

import (
	"errors"

	"github.com/tusklang/tusk/tokenizer"
)

type VarDecl struct {
	Name string
	Type *ASTNode
}

func (vd *VarDecl) Parse(lex []tokenizer.Token, i *int) error {

	*i++

	if lex[*i].Type != "varname" {
		return errors.New("expected a variable name")
	}

	vd.Name = lex[*i].Name

	//has a specified type
	if lex[*i+1].Name == ":" {
		*i += 2
		t, e := groupsToAST(groupSpecific(lex, 1, i))
		vd.Type = t[0]
		if e != nil {
			return e
		}
	}

	*i-- //the outer loop will incremenet for us

	return nil
}
