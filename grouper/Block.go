package grouper

import (
	"errors"

	"github.com/tusklang/tusk/tokenizer"
)

type Block struct {
	BlockType string
	Sub       []Group
}

var bmatches = map[string]string{
	"{": "}",
	"(": ")",
}

func (b *Block) Parse(lex []tokenizer.Token, i *int) error {

	if lex[*i].Type != "(" && lex[*i].Type != "{" {
		return errors.New("given lex is not a group")
	}

	b.BlockType = lex[*i].Type

	gcontent := Grouper(braceMatcher(lex, i, lex[*i].Type, bmatches[lex[*i].Type], true, ""))
	b.Sub = gcontent

	return nil
}
