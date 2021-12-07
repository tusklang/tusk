package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type This struct {
	tok tokenizer.Token
}

func (t *This) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {
	t.tok = lex[*i]
	return nil
}

func (t *This) GetMTok() tokenizer.Token {
	return t.tok
}

func (t *This) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	if !function.IsMethod {
		//error
		//cannot use `this` on non-methods
	}

	return data.NewVariable(function.LLFunc.Params[0], data.NewInstance(function.MethodClass))
}
