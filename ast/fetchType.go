package ast

import (
	"github.com/llir/llvm/ir/types"
)

//convert a group to a llvm type

func (compiler *Compiler) fetchValidType(name string) (v types.Type, e error) {

	var exists bool
	if v, exists = compiler.ValidTypes[name]; !exists {
		return nil, nil //error
	}

	return
}

func (compiler *Compiler) FetchType(class *types.StructType, g Group) (t types.Type, e error) {
	switch gt := g.(type) {
	case *DataType:
		return compiler.fetchValidType(gt.Type.Name)
	case *Function:

		if gt.Body != nil {
			return nil, nil //a function type must only be a function header, not including a body
		}

		var rt types.Type = &types.VoidType{}

		if gt.RetType != nil {
			rt, e = compiler.FetchType(class, gt.RetType.Group)

			if e != nil {
				return nil, e
			}
		}

		var params = make([]types.Type, len(gt.Params))

		for k, v := range gt.Params {
			params[k], e = compiler.FetchType(class, v.Type.Group)

			if e != nil {
				return nil, e
			}
		}

		return types.NewPointer(types.NewFunc(rt, params...)), nil
	case *VarRef:
		return compiler.fetchValidType(gt.Name)
	}

	//default case:

	//error
	return nil, nil //temp
}
