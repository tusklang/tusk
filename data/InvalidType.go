package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type InvalidType struct{}

func NewInvalidType() *InvalidType {
	return &InvalidType{}
}

func (n *InvalidType) LLVal(function *Function) value.Value {
	return nil
}

func (n *InvalidType) TType() Type {
	return n
}

func (n *InvalidType) Type() types.Type {
	return nil
}

func (n *InvalidType) TypeData() *TypeData {
	return NewTypeData("invalid")
}

func (n *InvalidType) InstanceV() value.Value {
	return nil
}

func (n *InvalidType) Alignment() uint64 {
	return 0
}

func (n *InvalidType) Default() constant.Constant {
	return nil
}

func (n *InvalidType) Equals(t Type) bool {
	return false
}
