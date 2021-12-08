package ast

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type Return struct {
	Val *ASTNode

	kwtok, rettok tokenizer.Token
}

func (r *Return) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {

	r.kwtok = lex[*i]

	*i++

	r.rettok = lex[*i]

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

	*i--

	return e
}

func (r *Return) GetMTok() tokenizer.Token {
	return r.kwtok
}

func (r *Return) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	var crvalLL value.Value

	if r.Val == nil {
		function.ActiveBlock.NewRet(nil)
		return nil
	}

	crval := r.Val.Group.Compile(compiler, class, r.Val, function) //compile the return val

	if !crval.TType().Equals(function.RetType()) {
		ocrv := crval
		crval = compiler.CastStore.RunCast(true, function.RetType(), crval, r.Val.Group, compiler, function, class)
		if crval == nil {
			//error
			//return type doesn't match the type here
			compiler.AddError(errhandle.NewCompileErrorFTok(
				"wrong return type",
				fmt.Sprintf("expected type %s to return but got %s", function.RetType().TypeData(), ocrv.TypeData()),
				r.rettok,
			))
			return nil
		}
	}

	crvalLL = crval.LLVal(function)
	function.ActiveBlock.NewRet(crvalLL)

	//all code after a return is unreachable, so it's placed into a dummy block
	//this block isn't ever compiled to llvm
	//it's just here to absorb all obsolete instructions
	function.ActiveBlock = ir.NewBlock("")

	return nil
}
