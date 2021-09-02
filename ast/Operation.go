package ast

import "github.com/tusklang/tusk/tokenizer"

type Operation struct {
	OpType string
	Token  *tokenizer.Token
}

func (o *Operation) Parse(lex []tokenizer.Token, i *int) error {

	o.Token = &lex[*i]
	o.OpType = lex[*i].Name

	return nil
}

func (o *Operation) Compile(compiler *Compiler, node *ASTNode, lvl int) {}
