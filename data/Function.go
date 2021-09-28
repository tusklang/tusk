package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Function struct {
	llfunc *ir.Func
	ret    Type
}

func NewFunc(f *ir.Func, ret Type) *Function {
	return &Function{
		llfunc: f,
		ret:    ret,
	}
}

func (f *Function) LLVal(block *ir.Block) value.Value {
	return f.llfunc
}

func (f *Function) RetType() Type {
	return f.ret
}

func (f *Function) Default() constant.Constant {
	return constant.NewNull(f.llfunc.Typ)
}

func (f *Function) TType() Type {
	return f
}

func (f *Function) Type() types.Type {
	return f.llfunc.Type()
}

func (f *Function) TypeString() string {
	return "func"
}
