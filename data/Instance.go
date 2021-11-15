package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Instance struct {
	Class *Class
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		Class: class,
	}
}

func (i *Instance) LLVal(block *ir.Block) value.Value {
	return nil
}

func (i *Instance) TType() Type {
	return i
}

func (i *Instance) Type() types.Type {
	return i.Class.Type()
}

func (i *Instance) TypeData() *TypeData {

	td := NewTypeData(i.Class.Name)
	td.AddFlag("instance")

	return td
}

func (i *Instance) InstanceV() value.Value {
	return nil
}

func (i *Instance) Equals(t Type) bool {
	return i.Class.Equals(t)
}

func (i *Instance) Default() constant.Constant {
	return constant.NewNull(types.NewPointer(i.Class.SType))
}

func (i *Instance) TypeSize() uint64 {
	return 8
}
