package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type This struct{}

func (t *This) Parse(lex []tokenizer.Token, i *int, stopAt []string) error {
	return nil
}

func (t *This) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	if !function.IsMethod {
		//error
		//cannot use `this` on non-methods
	}

	return data.NewInstVariable(function.LLFunc.Params[0], data.NewInstance(function.MethodClass))
}
