package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Function struct {
	llfunc *ir.Func
}

func NewFunc(f *ir.Func) *Function {
	return &Function{
		llfunc: f,
	}
}

func (f *Function) LLVal(block *ir.Block) value.Value {
	return f.llfunc
}

func (f *Function) Type() types.Type {
	return f.llfunc.Type()
}

func FuncTypeDefault(t types.Type) string {
	return "func"
}

func (f *Function) TypeString() string {
	return FuncTypeDefault(nil)
}
