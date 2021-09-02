package ast

import "github.com/tusklang/tusk/tokenizer"

type VarRef struct {
	Name string
}

func (vr *VarRef) Parse(lex []tokenizer.Token, i *int) error {
	vr.Name = lex[*i].Name
	return nil
}

func (vr *VarRef) Compile(compiler *Compiler, node *ASTNode, lvl int) {}
