package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Static struct {
	Declaration *ASTNode

	tok tokenizer.Token
}

func (s *Static) Parse(lex []tokenizer.Token, i *int, stopAt []string) (e error) {
	s.tok = lex[*i]
	return parseAccessSpec(s, lex, i)
}

func (s *Static) SetDecl(node *ASTNode) {
	s.Declaration = node
}

func (s *Static) GetMTok() tokenizer.Token {
	return s.tok
}

//cannot be compiled
func (s *Static) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return nil
}
