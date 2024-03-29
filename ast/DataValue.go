package ast

import (
	"strconv"

	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

type DataValue struct {
	Value data.Value
	tok   tokenizer.Token
}

func parseInt(n string) int64 {
	nv, _ := strconv.ParseInt(n, 10, 64)
	return nv
}

func parseFloat(n string) float64 {
	nv, _ := strconv.ParseFloat(n, 64)
	return nv
}

func parseBool(n string) bool {
	if n == "true" {
		return true
	} else {
		return false
	}
}

func getValue(tok tokenizer.Token) data.Value {
	switch tok.Type {
	case "int":
		return data.NewUntypedInteger(parseInt(tok.Name))
	case "float":
		return data.NewUntypedFloat(parseFloat(tok.Name))
	case "bool":
		return data.NewBoolean(parseBool(tok.Name))
	}

	return nil
}

func (dv *DataValue) Parse(lex []tokenizer.Token, i *int, stopAt []string) *errhandle.TuskError {
	dv.Value = getValue(lex[*i])
	dv.tok = lex[*i]
	return nil
}

func (dv *DataValue) GetMTok() tokenizer.Token {
	return dv.tok
}

func (dv *DataValue) Compile(compiler *Compiler, class *data.Class, node *ASTNode, function *data.Function) data.Value {
	return dv.Value
}
