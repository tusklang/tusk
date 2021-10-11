package data

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func LLTypToTusk(llt types.Type) Type {
	switch typ := llt.(type) {
	case *types.FuncType:

		var params = make([]*ir.Param, len(typ.Params))

		for k, v := range typ.Params {
			params[k] = ir.NewParam("", v)
		}

		var tfunc = NewFunc(ir.NewFunc(typ.Name(), typ.RetType, params...), LLTypToTusk(typ.RetType))

		return tfunc
	case *types.PointerType:

		switch typ.ElemType.(type) {
		case *types.FuncType:
			//a pointer function is regarded as a regular function in tusk
			return LLTypToTusk(typ.ElemType)
		}

		return NewPointer(LLTypToTusk(typ.ElemType))
	default:
		return NewPrimitive(llt)
	}
}
