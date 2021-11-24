package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type SliceArray struct {
	atype  Type
	decl   value.Value
	curlen value.Value
}

func NewSliceArray(atype Type, decl, curlen value.Value) *SliceArray {
	return &SliceArray{
		atype:  atype,
		decl:   decl,
		curlen: curlen,
	}
}

func (a *SliceArray) ValType() Type {
	return a.atype
}

func (a *SliceArray) LLVal(block *ir.Block) value.Value {
	return block.NewLoad(types.NewPointer(a.atype.Type()), a.decl)
}

func (a *SliceArray) TValue() Value {
	return a
}

func (a *SliceArray) TType() Type {
	return a
}

func (a *SliceArray) Type() types.Type {
	return types.NewPointer(a.atype.Type())
}

func (a *SliceArray) TypeData() *TypeData {
	td := NewTypeData("array")
	td.AddFlag("slice")
	return td
}

func (a *SliceArray) InstanceV() value.Value {
	return nil
}

func (a *SliceArray) Default() constant.Constant {
	return constant.NewNull(types.NewPointer(a.atype.Type()))
}

func (a *SliceArray) Equals(Type) bool {
	return false
}

func (a *SliceArray) Alignment() uint64 {
	return 8
}
