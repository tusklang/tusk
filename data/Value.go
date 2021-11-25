package data

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Value interface {
	LLVal(function *Function) value.Value
	TType() Type
	Type() types.Type
	TypeData() *TypeData
	InstanceV() value.Value
}
