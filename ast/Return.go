package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Return struct {
	Val *ASTNode
}

func (r *Return) Parse(lex []tokenizer.Token, i *int, stopAt []string) error {

	*i++

	retval := braceMatcher(lex, i, allopeners, allclosers, false, "terminator")
	retvalAST, e := groupsToAST(grouper(retval))

	r.Val = retvalAST[0]

	return e
}

func (r *Return) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	crval := r.Val.Group.Compile(compiler, class, r.Val, function) //compile the return val

	if !crval.TType().Equals(function.RetType()) {
		crval = compiler.CastStore.RunCast(true, function.RetType(), crval, compiler, function, class)
	}

	function.ActiveBlock.NewRet(crval.LLVal(function))
	return nil
}
