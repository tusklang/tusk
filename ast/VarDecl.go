package ast

import (
	"errors"

	"github.com/tusklang/tusk/tokenizer"
)

type VarDecl struct {
	Name string
}

func (vd *VarDecl) Parse(lex []tokenizer.Token, i *int) error {

	*i++

	if lex[*i].Type != "varname" {
		return errors.New("expected a variable name")
	}

	vd.Name = lex[*i].Name

	return nil
}
