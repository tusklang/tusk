package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Boolean struct {
	val int64
}

func NewBoolean(val bool) *Boolean {

	var ival int64 //booleans are stored as integers
	if val {
		ival = 1
	}

	return &Boolean{
		val: ival,
	}
}

func (b *Boolean) LLVal(function *Function) value.Value {
	return constant.NewInt(types.I1, b.val)
}

func (b *Boolean) TType() Type {
	return NewNamedPrimitive(types.I1, "bool")
}

func (b *Boolean) Type() types.Type {
	return types.I1
}

func (b *Boolean) TypeData() *TypeData {
	return NewTypeData("bool")
}

func (b *Boolean) InstanceV() value.Value {
	return nil
}
