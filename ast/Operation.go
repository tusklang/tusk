package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/tokenizer"
)

type Operation struct {
	OpType string
	Token  *tokenizer.Token
}

func (o *Operation) Parse(lex []tokenizer.Token, i *int) error {

	o.Token = &lex[*i]
	o.OpType = lex[*i].Name

	return nil
}

func (o *Operation) Compile(compiler *Compiler, class *types.StructType, node *ASTNode, block *ir.Block) value.Value {

	var (
		l = node.Left[0]
		r = node.Right[0]
	)

	lc := l.Group.Compile(compiler, class, l, block)
	rc := r.Group.Compile(compiler, class, r, block)

	opString := typeToString(lc.Type()) + " " + o.OpType + " " + typeToString(rc.Type())

	return compiler.Operations[opString](lc, rc, block)
}
