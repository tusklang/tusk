package ast

import (
	"errors"

	"github.com/llir/llvm/ir"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type Block struct {
	BlockType string
	Sub       []*ASTNode
}

var bmatches = map[string]string{
	"{": "}",
	"(": ")",
}

func (b *Block) Parse(lex []tokenizer.Token, i *int) (e error) {

	if lex[*i].Type != "(" && lex[*i].Type != "{" {
		return errors.New("given lex is not a group")
	}

	b.BlockType = lex[*i].Type

	gcontent := grouper(braceMatcher(lex, i, []string{lex[*i].Type}, []string{bmatches[lex[*i].Type]}, true, ""))

	b.Sub, e = groupsToAST(gcontent)

	if e != nil {
		return e
	}

	return nil
}

func (b *Block) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {

	if b.BlockType == "(" {
		return b.Sub[0].Group.Compile(compiler, class, b.Sub[0], block)
	} else if b.BlockType == "fncallb" {
		//if it's a function call block

		var args = data.NewFnCallBlock()

		for _, v := range b.Sub {
			d := v.Group.Compile(compiler, class, b.Sub[0], block)
			args.Args = append(args.Args, d)
		}

		return args
	}

	for _, v := range b.Sub {
		v.Group.Compile(compiler, class, v, block)
	}

	return nil
}
