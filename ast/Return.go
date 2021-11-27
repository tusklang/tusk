package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Return struct {
	Val *ASTNode

	tok tokenizer.Token
}

func (r *Return) Parse(lex []tokenizer.Token, i *int, stopAt []string) error {

	r.tok = lex[*i]

	*i++

	retval := braceMatcher(lex, i, allopeners, allclosers, false, "terminator")
	retvalAST, e := groupsToAST(grouper(retval))

	//if there is no retval, then just return the `e`
	if len(retvalAST) == 0 {
		return e
	}

	r.Val = retvalAST[0]

	return e
}

func (r *Return) GetMTok() tokenizer.Token {
	return r.tok
}

func (r *Return) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	if r.Val == nil {
		function.ActiveBlock.NewRet(nil)
		return nil
	}

	crval := r.Val.Group.Compile(compiler, class, r.Val, function) //compile the return val

	if !crval.TType().Equals(function.RetType()) {
		crval = compiler.CastStore.RunCast(true, function.RetType(), crval, compiler, function, class)
	}

	function.ActiveBlock.NewRet(crval.LLVal(function))
	return nil
}
