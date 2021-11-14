package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"

	"github.com/llir/llvm/ir/types"
)

type Integer struct {
	value *constant.Int
}

func NewInteger(i *constant.Int) *Integer {
	return &Integer{
		value: i,
	}
}

func (i *Integer) LLVal(block *ir.Block) value.Value {
	return i.value
}

func (i *Integer) TType() Type {
	return NewPrimitive(i.Type())
}

func (i *Integer) Type() types.Type {
	return i.value.Type()
}

func (i *Integer) TypeData() *TypeData {
	return NewTypeData(i.Type().LLString())
}

func (i *Integer) InstanceV() value.Value {
	return nil
}
