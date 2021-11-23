package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type FixedArray struct {
	atype  Type
	decl   value.Value
	curlen value.Value
	length uint64
}

func NewFixedArray(atype Type, decl, curlen value.Value, length uint64) *FixedArray {
	return &FixedArray{
		atype:  atype,
		decl:   decl,
		curlen: curlen,
		length: length,
	}
}

func (a *FixedArray) GetIndex(block *ir.Block, idx Value) Value {
	gept := types.NewArray(a.length, a.atype.Type())
	gep := block.NewGetElementPtr(gept, a.LLVal(block), constant.NewInt(types.I32, 0), idx.LLVal(block))
	gep.InBounds = true
	return NewVariable(
		gep,
		nil,
		a.atype,
	)
}

func (a *FixedArray) LLVal(block *ir.Block) value.Value {
	return a.decl
}

func (a *FixedArray) TValue() Value {
	return a
}

func (a *FixedArray) TType() Type {
	return a
}

func (a *FixedArray) Type() types.Type {
	return a.decl.Type()
}

func (a *FixedArray) TypeData() *TypeData {
	td := NewTypeData("array")
	td.AddFlag("fixed")
	return td
}

func (a *FixedArray) InstanceV() value.Value {
	return nil
}

func (a *FixedArray) Default() constant.Constant {
	return constant.NewUndef(a.Type())
}

func (a *FixedArray) Equals(Type) bool {
	return false
}

func (a *FixedArray) Alignment() uint64 {
	return 16
}
