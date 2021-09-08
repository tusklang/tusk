package ast

import (
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
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
		return data.NewInteger(constant.NewInt(types.I32, parseInt(tok.Name)))
	case "float":
		return data.NewFloat(constant.NewFloat(types.Float, parseFloat(tok.Name)))
	}

	return nil
}

func (dv *DataValue) Parse(lex []tokenizer.Token, i *int) error {
	dv.Value = getValue(lex[*i])
	return nil
}

func (dv *DataValue) Compile(compiler *Compiler, class *data.Class, node *ASTNode, block *ir.Block) data.Value {
	return dv.Value
}
