package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type IfStatement struct {
	Condition []*ASTNode
	Body      []*ASTNode
	ElseBody  []*ASTNode
}

func (is *IfStatement) Parse(lex []tokenizer.Token, i *int) error {
	return ifwhileParse(is, lex, i)
}

func (is *IfStatement) SetCond(g []*ASTNode) {
	is.Condition = g
}

func (is *IfStatement) SetBody(g []*ASTNode) {
	is.Body = g
}

func (is *IfStatement) Type() string {
	return "if"
}

func (is *IfStatement) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	cond := is.Condition[0].Group.Compile(compiler, class, is.Condition[0], function)

	var (
		//create a true and false block
		trueblock  = function.LLFunc.NewBlock("")
		falseblock = function.LLFunc.NewBlock("")
	)

	function.ActiveBlock.NewCondBr(cond.LLVal(function.ActiveBlock), trueblock, falseblock)

	//block to store the intructions that come after the if statement
	var restActs = function.LLFunc.NewBlock("")

	gotoRestActs := ir.NewBr(restActs)
	function.PushTermStack(
		gotoRestActs,
	)

	function.ActiveBlock = trueblock //temporarily change the active block to the true block, so the compiler appends instructions to it
	is.Body[0].Group.Compile(compiler, class, is.Body[0], function)

	if is.ElseBody != nil {
		//if the else body isn't empty
		function.ActiveBlock = falseblock
		is.ElseBody[0].Group.Compile(compiler, class, is.ElseBody[0], function)
	}

	//if the latest terminator was not used, still pop it
	if gotoRestActs == function.LastTermStack() {
		function.PopTermStack()
	}

	//we change the active block to another block that stores all the further instructions
	function.ActiveBlock = restActs

	if val := function.PopTermStack(); val != nil {
		function.ActiveBlock.Term = val
	}

	//at the end of the ifs, go to the next instructions block, defined above
	//only if there isn't a terminator already ^
	if trueblock.Term == nil {
		trueblock.NewBr(function.ActiveBlock)
	}
	if falseblock.Term == nil {
		falseblock.NewBr(function.ActiveBlock)
	}

	return nil
}
