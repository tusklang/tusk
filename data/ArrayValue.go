package data

import "github.com/llir/llvm/ir"

type ArrayValue interface {
	GetIndex(*ir.Block, Value) Value
}
