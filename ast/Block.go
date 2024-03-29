package ast

import (
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type Block struct {
	BlockType string
	Sub       []*ASTNode

	tok tokenizer.Token
}

var bmatches = map[string]string{
	"{": "}",
	"(": ")",
}

var allopeners = []string{"{", "("}
var allclosers = []string{"}", ")"}

func (b *Block) Parse(lex []tokenizer.Token, i *int, stopAt []string) (e *errhandle.TuskError) {

	b.tok = lex[*i]

	b.BlockType = lex[*i].Type

	contbm, e := braceMatcher(lex, i, []string{lex[*i].Type}, []string{bmatches[lex[*i].Type]}, true, "")

	if e != nil {
		return e
	}

	gcontent, e := grouper(contbm)

	if e != nil {
		return e
	}

	b.Sub, e = groupsToAST(gcontent)

	if e != nil {
		return e
	}

	return nil
}

func (b *Block) GetMTok() tokenizer.Token {
	return b.tok
}

func (b *Block) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {

	if b.BlockType == "(" {
		return b.Sub[0].Group.Compile(compiler, class, b.Sub[0], function)
	} else if b.BlockType == "fncallb" {
		//if it's a function call block

		var args = data.NewFnCallBlock()

		for _, v := range b.Sub {
			d := v.Group.Compile(compiler, class, v, function)
			args.Args = append(args.Args, d)
		}

		return args
	}

	for _, v := range b.Sub {
		v.Group.Compile(compiler, class, v, function)
	}

	return nil
}
