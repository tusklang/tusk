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

func espAppend(b []byte, val byte) []byte {
	switch val {
	case 'n':
		return append(b, '\n')
	case '0':
		return append(b, 0)
	default:
		return append(b, val)
	}
}

func escString(b []byte) []byte {
	//escape a string

	var escaped bool
	var fin []byte

	for _, v := range b {
		if escaped {
			fin = espAppend(fin, v)
			escaped = false
			continue
		}

		if v == '\\' {
			escaped = true
			continue
		}

		fin = append(fin, v)
	}

	return fin
}

func (s *String) Parse(lex []tokenizer.Token, i *int) error {

	sv := lex[*i].Name

	sv = strings.TrimSuffix(strings.TrimPrefix(sv, "\""), "\"") //remove the leading and trailing quotes

	s.dstring = data.NewString(
		escString(append([]byte(sv), 0)),
	)

	return nil
}

func (s *String) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	sd := compiler.Module.NewGlobalDef("", constant.NewCharArray(s.dstring.CharArray))
	s.dstring.Init(sd)
	return s.dstring
}
