package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type VariedLengthArray struct {
	atype  Type
	decl   value.Value
	curlen value.Value
	length value.Value
}

func NewVariedLengthArray(atype Type, decl, curlen value.Value, length value.Value) *VariedLengthArray {
	return &VariedLengthArray{
		atype:  atype,
		decl:   decl,
		curlen: curlen,
		length: length,
	}
}

func (a *VariedLengthArray) ValType() Type {
	return a.atype
}

func (a *VariedLengthArray) LLVal(block *ir.Block) value.Value {
	return a.decl
}

func (a *VariedLengthArray) TValue() Value {
	return a
}

func (a *VariedLengthArray) TType() Type {
	return a
}

func (a *VariedLengthArray) Type() types.Type {
	return a.decl.Type()
}

func (a *VariedLengthArray) TypeData() *TypeData {
	td := NewTypeData("array")
	td.AddFlag("varied")
	return td
}

func (a *VariedLengthArray) InstanceV() value.Value {
	return nil
}

func (a *VariedLengthArray) Default() constant.Constant {
	return constant.NewUndef(a.Type())
}

func (a *VariedLengthArray) Equals(t Type) bool {
	switch c := t.(type) {
	case *VariedLengthArray:
		return a.ValType().Equals(c.ValType())
	}
	return false
}

func (a *VariedLengthArray) Alignment() uint64 {
	return 16
}
