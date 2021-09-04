package ast

import (
	"github.com/llir/llvm/ir/types"
)

//convert a group to a llvm type

func fetchBasicType(name string) (types.Type, error) {

	switch name {
	case "i64":
		return types.I64, nil
	case "i32":
		return types.I32, nil
	case "i16":
		return types.I16, nil
	case "i8":
		return types.I8, nil
	}

	//default case:

	//error
	return nil, nil //implement later
}

func fetchType(g Group) (t types.Type, e error) {
	switch gt := g.(type) {
	case *DataType:
		return fetchBasicType(gt.Type.Name)
	case *Function:

		if gt.Body != nil {
			return nil, nil //a function type must only be a function header, not including a body
		}

		var rt types.Type = &types.VoidType{}

		if gt.RetType != nil {
			rt, e = fetchType(gt.RetType.Group)

			if e != nil {
				return nil, e
			}
		}

		var params = make([]types.Type, len(gt.Params))

		for k, v := range gt.Params {
			params[k], e = fetchType(v.Type.Group)

			if e != nil {
				return nil, e
			}
		}

		return types.NewPointer(types.NewFunc(rt, params...)), nil
	}

	//default case:

	//error
	return nil, nil //temp
}
