package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type Return struct {
	Val *ASTNode

	tok tokenizer.Token
}

func (r *Return) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {

	r.tok = lex[*i]

	*i++

	retval, e := braceMatcher(lex, i, allopeners, allclosers, false, "terminator")

	if e != nil {
		return e
	}

	rg, e := grouper(retval)

	if e != nil {
		return e
	}

	retvalAST, e := groupsToAST(rg)

	if e != nil || len(retvalAST) == 0 /*if there isn't a retval, don't go any further*/ {
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
		crval = compiler.CastStore.RunCast(true, function.RetType(), crval, r.Val.Group, compiler, function, class)
	}

	function.ActiveBlock.NewRet(crval.LLVal(function))
	return nil
}
