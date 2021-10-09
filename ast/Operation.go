package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Operation struct {
	OpType string
	Token  *tokenizer.Token
}

func (o *Operation) Parse(lex []tokenizer.Token, i *int) error {

	o.Token = &lex[*i]
	o.OpType = lex[*i].Name

	return nil
}

func (o *Operation) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	var (
		lc data.Value
		rc data.Value
	)

	//parse the left and right operands
	//only if they exist
	//some operators are only one side (e.g. !, ~, etc...)
	if len(node.Left) != 0 {
		lc = node.Left[0].Group.Compile(compiler, class, node.Left[0], function)
	}
	if len(node.Right) != 0 {
		rc = node.Right[0].Group.Compile(compiler, class, node.Right[0], function)
	}

	if o.OpType == "." {

		//if it's a . operator (accessing variables within a class/instance)
		//then the right operand must be an undeclared variable (name of the field to access)
		//it can be a variable if the name of the field is equal to a variable's name, so the compiler fetches the variable, instead of reporting an udvar
		switch rc.(type) {
		case *data.Variable:
			rc = data.NewUndeclaredVar(node.Right[0].Group.(*VarRef).Name)
		}
	}

	rop := compiler.OperationStore.RunOperation(lc, rc, o.OpType, compiler, function.ActiveBlock)
	return rop
}
