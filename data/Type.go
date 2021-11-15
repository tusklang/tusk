package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Type interface {
	TType() Type
	Type() types.Type
	TypeData() *TypeData
	Default() constant.Constant
	Equals(Type) bool
	InstanceV() value.Value
	TypeSize() uint64
}
