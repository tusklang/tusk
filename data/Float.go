package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"

	"github.com/llir/llvm/ir/types"
)

type Float struct {
	value *constant.Float
}

func NewFloat(f *constant.Float) *Float {
	return &Float{
		value: f,
	}
}

func (f *Float) LLVal(block *ir.Block) value.Value {
	return f.value
}

func (f *Float) TType() Type {
	return NewPrimitive(f.Type())
}

func (f *Float) Type() types.Type {
	return f.value.Type()
}

func (f *Float) TypeData() *TypeData {
	return NewTypeData(f.value.Type().LLString())
}
