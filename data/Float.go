package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"

	"github.com/llir/llvm/ir/types"
)

type Float struct {
	value *constant.Float

	//untyped
	untyped bool
	UTypVal float64
}

func NewFloat(f *constant.Float) *Float {
	return &Float{
		value: f,
	}
}

func NewUntypedFloat(v float64) *Float {
	return &Float{
		untyped: true,
		UTypVal: v,
	}
}

func (f *Float) LLVal(block *ir.Block) value.Value {
	return f.value
}

func (f *Float) TType() Type {

	if f.untyped {
		return NewUntypeFloatType(f)
	}

	return NewPrimitive(f.Type())
}

func (f *Float) Type() types.Type {

	if f.untyped {
		return nil
	}

	return f.value.Type()
}

func (f *Float) TypeData() *TypeData {
	return f.TType().TypeData()
}

func (f *Float) InstanceV() value.Value {
	return nil
}
