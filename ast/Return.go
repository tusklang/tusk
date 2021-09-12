package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Return struct {
	Val *ASTNode
}

func (r *Return) Parse(lex []tokenizer.Token, i *int) error {

	*i++

	retval := braceMatcher(lex, i, allopeners, allclosers, false, "terminator")
	retvalAST, e := groupsToAST(grouper(retval))

	r.Val = retvalAST[0]

	return e
}

func (r *Return) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {
	crval := r.Val.Group.Compile(compiler, class, r.Val, block) //compile the return val
	block.NewRet(crval.LLVal(block))
	return nil
}
