package ast

import (
	"strings"

	"github.com/llir/llvm/ir/constant"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type String struct {
	dstring *data.String
}

func (s *String) Parse(lex []tokenizer.Token, i *int) error {

	sv := lex[*i].Name

	sv = strings.TrimSuffix(strings.TrimPrefix(sv, "\""), "\"") //remove the leading and trailing quotes

	s.dstring = data.NewString([]byte(sv))

	return nil
}

func (s *String) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	sd := compiler.Module.NewGlobalDef("", constant.NewCharArray(s.dstring.CharArray))
	s.dstring.Init(sd, compiler.NewString)
	return s.dstring
}
