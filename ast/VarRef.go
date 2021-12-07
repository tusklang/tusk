package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type VarRef struct {
	Name string
	tok  tokenizer.Token
}

func (vr *VarRef) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {
	vr.Name = lex[*i].Name
	vr.tok = lex[*i]
	return nil
}

func (vr *VarRef) GetMTok() tokenizer.Token {
	return vr.tok
}

func (vr *VarRef) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	fetched := compiler.FetchVar(vr.Name)

	var global, pure bool

	if fetched == nil {

		//check the class' static variables if there is no variable declared with x name

		_fetched := class.Static[vr.Name]

		if _fetched != nil {
			fetched = _fetched.Value
			global = true
			pure = _fetched.Pure
		}

		if fetched == nil {
			//if there still isn't a variable with that name, it's an "undeclared variable"
			return data.NewUndeclaredVar(vr.Name)
		}

	}

	if function != nil && function.IsPure && global && !pure {
		//if the function we're in is pure
		//make sure that we're not using any globals, or global functions
		//(other functions that aren't pure***)
		compiler.AddError(errhandle.NewCompileErrorFTok(
			"pure function accessing impure global",
			"cannot access globals from pure functions",
			vr.tok,
		))
		return data.NewInvalidType()
	}

	return fetched
}
