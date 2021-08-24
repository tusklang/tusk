package grouper

import "github.com/tusklang/tusk/tokenizer"

type Default struct {
	Token tokenizer.Token
}

func (d *Default) Parse(lex []tokenizer.Token, i *int) error {
	d.Token = lex[*i]
	return nil
}
