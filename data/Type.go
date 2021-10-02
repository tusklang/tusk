package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type Type interface {
	TType() Type
	Type() types.Type
	TypeData() *TypeData
	Default() constant.Constant
}
