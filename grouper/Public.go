package grouper

import (
	"github.com/tusklang/tusk/tokenizer"
)

type Public struct {
	Declaration Group
}

func (p *Public) Parse(lex []tokenizer.Token, i *int) error {
	*i++
	g := groupSpecific(lex, 1, i)
	p.Declaration = g[0]
	*i-- //decrement because the outer grouper function already increments once
	return nil
}
