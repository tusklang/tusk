package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type NullValue struct {
	tok tokenizer.Token
}

func (nv *NullValue) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {
	nv.tok = lex[*i]
	return nil
}

func (nv *NullValue) GetMTok() tokenizer.Token {
	return nv.tok
}

func (dv *NullValue) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return data.NewNull()
}
