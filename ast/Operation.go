package ast

import (
	"fmt"

	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type Operation struct {
	OpType string
	tok    tokenizer.Token
}

func (o *Operation) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {

	o.tok = lex[*i]
	o.OpType = lex[*i].Name

	return nil
}

func (o *Operation) GetMTok() tokenizer.Token {
	return o.tok
}

func (o *Operation) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	//parse the left and right operands
	//only if they exist
	//some operators are only one side (e.g. !, ~, etc...)
	var lc, rc data.Value
	var lcg, rcg Group
	if len(node.Left) != 0 {
		lcg = node.Left[0].Group
		lc = lcg.Compile(compiler, class, node.Left[0], function)
	}
	if len(node.Right) != 0 {
		rcg = node.Right[0].Group
		rc = rcg.Compile(compiler, class, node.Right[0], function)
	}

	if o.OpType == "." {

		//if it's a . operator (accessing variables within a class/instance)
		//then the right operand must be an undeclared variable (name of the field to access)
		//it can be a variable if the name of the field is equal to a variable's name, so the compiler fetches the variable, instead of reporting an udvar
		switch rc.(type) {
		case *data.Variable:
			rc = data.NewUndeclaredVar(node.Right[0].Group.(*VarRef).Name)
		case *data.Function:
			rc = data.NewUndeclaredVar(rc.(*data.Function).GetLName())
		case *data.Class:
			rc = data.NewUndeclaredVar(rc.TypeData().Name())
		}

	}

	rop := compiler.OperationStore.RunOperation(lc, rc, lcg, rcg, o.OpType, compiler, function, class)

	if rop == nil {

		var lct, rct string

		if lc == nil {
			lct = "none"
		} else {
			lct = lc.TType().TypeData().String()
		}

		if rc == nil {
			rct = "none"
		} else {
			rct = rc.TType().TypeData().String()
		}

		compiler.AddError(errhandle.NewCompileErrorFTok(
			"invalid operation",
			fmt.Sprintf("'%s' operation not found for %s and %s", o.OpType, lct, rct),
			o.GetMTok(),
		))
		return data.NewInvalidType()
	}

	return rop
}
