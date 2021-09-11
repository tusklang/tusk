package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/data"
)

//intialize string utility things
//(string type, and new string func)
func initString(module *ir.Module) (*data.Class, *ir.Func) {

	stype := types.NewStruct() //create a new structure (representing a class)

	stype.Fields = []types.Type{types.I8Ptr, types.I32}

	module.NewTypeDef("tusk.string", stype) //create the typedef in llvm

	sptr := ir.NewParam("sptr", types.I8Ptr)
	slen := ir.NewParam("slen", types.I32)

	stringf := module.NewFunc("tusk.newstring", stype,
		sptr,
		slen,
	)
	fbod := stringf.NewBlock("")

	salc := fbod.NewAlloca(stype)                                                                                //allocate a new string type
	salcSPtr := fbod.NewGetElementPtr(stype, salc, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)) //get the first element in the salc struct (sptr)
	fbod.NewStore(sptr, salcSPtr)                                                                                //store the given string in the salcSPtr
	salcSLen := fbod.NewGetElementPtr(stype, salc, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1)) //get the second element in the salc struct (slen)
	fbod.NewStore(slen, salcSLen)                                                                                //store the given string in the salcSLen
	loadSalc := fbod.NewLoad(stype, salc)                                                                        //load the salc variable to remove the reference
	fbod.NewRet(loadSalc)                                                                                        //return the allocated string

	return data.NewClass("tusk.string", stype, nil), stringf
}
