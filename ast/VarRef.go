package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type VarRef struct {
	Name string
}

func (vr *VarRef) Parse(lex []tokenizer.Token, i *int) error {
	vr.Name = lex[*i].Name
	return nil
}

func (vr *VarRef) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {
	fetched := compiler.FetchVar(vr.Name)

	if fetched == nil {
		return data.NewUndeclaredVar(vr.Name)
	}

	//it's an un-referenceable variable
	//(mostly used for types as variables)
	//so we just return the value of it, instead of the pointer
	if fetched.UnReferenceable {
		return fetched.FetchVal()
	}

	return fetched
}
