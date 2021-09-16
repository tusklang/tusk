package ast

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Construct struct {
	FnObj *Function
}

func (c *Construct) Parse(lex []tokenizer.Token, i *int) error {

	var fnobj = &Function{}
	e := fnobj.Parse(lex, i) //functions and constructors are (surprisingly enough :p) structured the same

	if e != nil { //if the function parse returned an error
		return e
	}

	if fnobj.RetType != nil { //constructors cannot return anything
		return errors.New("constructor cannot include a return type")
	}

	c.FnObj = fnobj

	return nil
}

//cannot be compiled like this
func (c *Construct) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {
	return nil
}

func (c *Construct) CompileConstructor(compiler *Compiler, class *data.Class, block *ir.Block, initval value.Value) error {

	var params = make([]*ir.Param, len(c.FnObj.Params))

	for k, v := range c.FnObj.Params {
		params[k] = ir.NewParam(
			v.Name,
			v.Type.Group.Compile(compiler, class, v.Type, block).Type(),
		)
	}

	//alter the params of the original init func
	block.Parent.Params = params

	//compile the constructor into a function
	constructor := c.FnObj.Compile(compiler, class, nil, block)

	//convert the params into args to call the new llvm ir func ^
	var args = make([]value.Value, len(c.FnObj.Params))

	for k, v := range params {
		args[k] = v
	}

	block.NewCall(constructor.LLVal(block), args...)

	return nil
}
