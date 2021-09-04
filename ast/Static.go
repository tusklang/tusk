package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/tokenizer"
)

type Static struct {
	Declaration *ASTNode
}

func (s *Static) Parse(lex []tokenizer.Token, i *int) (e error) {
	return parseAccessSpec(s, lex, i)
}

func (s *Static) SetDecl(node *ASTNode) {
	s.Declaration = node
}

//cannot be compiled
func (s *Static) Compile(class *types.StructType, node *ASTNode) {}
