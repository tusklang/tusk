package ast

import (
	"strconv"

	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/tokenizer"
)

type DataValue struct {
	Value data.Value
}

func parseInt(n string) int64 {
	nv, _ := strconv.ParseInt(n, 10, 64)
	return nv
}

func parseFloat(n string) float64 {
	nv, _ := strconv.ParseFloat(n, 64)
	return nv
}

func getValue(tok tokenizer.Token) data.Value {
	switch tok.Type {
	case "int":
		return data.NewUntypedInteger(parseInt(tok.Name))
	case "float":
		return data.NewUntypedFloat(parseFloat(tok.Name))
	}

	return nil
}

func (dv *DataValue) Parse(lex []tokenizer.Token, i *int) error {
	dv.Value = getValue(lex[*i])
	return nil
}

func (dv *DataValue) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return dv.Value
}
