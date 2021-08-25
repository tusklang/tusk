package grouper

import "github.com/tusklang/tusk/tokenizer"

type DataType struct {
	DType tokenizer.Token
}

func (dt *DataType) Parse(lex []tokenizer.Token, i *int) error {
	dt.DType = lex[*i]
	return nil
}
