package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Value interface {
	LLVal(block *ir.Block) value.Value
	Type() types.Type
	TypeString() string
}
