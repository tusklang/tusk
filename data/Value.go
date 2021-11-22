package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Value interface {
	LLVal(block *ir.Block) value.Value
	TValue() Value
	TType() Type
	Type() types.Type
	TypeData() *TypeData
	InstanceV() value.Value
}
