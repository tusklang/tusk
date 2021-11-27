package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type FnCallBlock struct {
	Args []Value
}

func NewFnCallBlock() *FnCallBlock {
	return &FnCallBlock{}
}

func (fcb *FnCallBlock) LLVal(function *Function) value.Value {
	return nil
}

func (fcb *FnCallBlock) TType() Type {
	return fcb
}

func (fcb *FnCallBlock) Type() types.Type {
	return nil
}

func (fcb *FnCallBlock) TypeData() *TypeData {
	return NewTypeData("fncallb")
}

func (fcb *FnCallBlock) InstanceV() value.Value {
	return nil
}

func (fcb *FnCallBlock) Alignment() uint64 {
	return 0
}

func (fcb *FnCallBlock) Default() constant.Constant {
	return nil
}

func (fcb *FnCallBlock) Equals(t Type) bool {
	return false
}
