package data

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

//get the default value for x type
func GetDefault(typ types.Type) constant.Constant {
	switch typ.(type) {
	case *types.IntType: //integers default to 0
		return constant.NewInt(types.I32, 0)
	default: //everything else defaults to null
		return &constant.Null{}
	}
}
