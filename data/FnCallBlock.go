package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type FnCallBlock struct {
	Args []Value
}

func NewFnCallBlock() *FnCallBlock {
	return &FnCallBlock{}
}

func (fcb *FnCallBlock) LLVal(block *ir.Block) value.Value {
	return nil
}

func (fcb *FnCallBlock) TValue() Value {
	return fcb
}

func (fcb *FnCallBlock) TType() Type {
	return nil
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
