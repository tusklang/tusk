package ast

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type Construct struct {
	FnObj *Function

	//during compilation
	cfnc   *data.Function
	params []*ir.Param

	tok tokenizer.Token
}

func (c *Construct) Parse(lex []tokenizer.Token, i *int, stopAt []string) error {

	c.tok = lex[*i]

	*i++

	if lex[*i].Type != "fn" {
		return errors.New("constructors must be functions")
	}

	var fnobj = &Function{}
	e := fnobj.Parse(lex, i, stopAt) //functions and constructors are (surprisingly enough :p) structured the same

	fnobj.isMethod = true

	if e != nil { //if the function parse returned an error
		return e
	}

	if fnobj.RetType != nil { //constructors cannot return anything
		return errors.New("constructor cannot include a return type")
	}

	c.FnObj = fnobj

	return nil
}

func (c *Construct) GetMTok() tokenizer.Token {
	return c.tok
}

//cannot be compiled like this
func (c *Construct) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}

func (c *Construct) CompileSig(compiler *Compiler, class *data.Class, function *data.Function) error {
	var params = make([]*ir.Param, len(c.FnObj.Params))
	var ptypes []data.Type

	for k, v := range c.FnObj.Params {
		typ := v.Type.Group.Compile(compiler, class, v.Type, function)
		params[k] = ir.NewParam(
			v.Name,
			typ.Type(),
		)
		ptypes = append(ptypes, typ.TType())
	}

	//alter the params of the original init func
	function.LLFunc.Params = params
	//add ptypes for tusk
	function.ParamTypes = ptypes

	c.cfnc = function
	c.params = params
	return nil
}

func (c *Construct) CompileConstructor(compiler *Compiler, class *data.Class) error {

	function := c.cfnc

	//compile the constructor into a function
	_constructor := c.FnObj.Compile(compiler, class, nil, function)

	//make sure the next decl is a function
	switch _constructor.(type) {
	case *data.Function:
	default:
		//error
		//expected a function
		compiler.AddError(errhandle.NewTuskErrorFTok(
			"constructors must be functions",
			"",
			c.tok,
		))
		return nil
	}

	constructor := _constructor.(*data.Function)

	if !constructor.RetType().Equals(data.NewPrimitive(types.Void)) {
		//error
		//constructors cannot have return types
		compiler.AddError(errhandle.NewTuskErrorFTok(
			"constructors cannot have return types",
			"",
			c.tok,
		))
		return nil
	}

	//convert the params into args to call the new llvm ir func ^
	var args = make([]value.Value, len(c.FnObj.Params))

	for k, v := range c.params {
		args[k] = v
	}

	function.ActiveBlock.NewCall(constructor.LLVal(function), append([]value.Value{class.ConstructAlloc}, args...)...)

	return nil
}
