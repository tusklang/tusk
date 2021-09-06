package ast

import (
	"strconv"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/tusklang/tusk/tokenizer"
)

//used for DataValue action

func parseInt(n string) int64 {
	nv, _ := strconv.ParseInt(n, 10, 64)
	return nv
}

func parseFloat(n string) float64 {
	nv, _ := strconv.ParseFloat(n, 64)
	return nv
}

func getValue(tok tokenizer.Token) value.Value {
	switch tok.Type {
	case "int":
		return constant.NewInt(types.I32, parseInt(tok.Name))
	case "float":
		return constant.NewFloat(types.Float, parseFloat(tok.Name))
	}

	return nil
}
