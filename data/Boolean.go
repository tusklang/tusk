package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Boolean struct {
	val *Integer
}

func NewBoolean(val bool) *Boolean {

	var ival int64 //booleans are stored as integers
	if val {
		ival = 1
	}

	return &Boolean{
		val: NewInteger(constant.NewInt(types.I1, ival)),
	}
}

func (b *Boolean) LLVal(block *ir.Block) value.Value {
	return b.val.LLVal(block)
}

func (b *Boolean) TType() Type {
	return NewPrimitive(types.I1)
}

func (b *Boolean) Type() types.Type {
	return types.I32
}

func (b *Boolean) TypeData() *TypeData {
	return NewTypeData("bool")
}

func (b *Boolean) InstanceV() value.Value {
	return nil
}
