package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type NullValue struct{}

func (nv *NullValue) Parse(lex []tokenizer.Token, i *int, stopAt []string) error {
	return nil
}

func (dv *NullValue) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return data.NewNull()
}
