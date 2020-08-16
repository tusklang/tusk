//The omm package allows go programs to call omm programs
package omm

import "lang/interpreter"
import "lang/types"
import "oat/encoding"

//Create a new instance of an omm script from an oat file
func InstanceFromOat(filename string) (*types.Instance, error) {
	oatv, e := oatenc.OatDecode(filename, 0)
	if e != nil {
		return nil, e
	}

	var ins *types.Instance
	
	for k, v := range oatv {
		ins.Allocate(k, interpreter.Interpreter(ins, v, []string{ "at goat binder" }).Exp)
	}

	return ins, nil
}

//Call a (global) function given an instance
func CallFunc(ins *types.Instance, funcname string, args ...*types.OmmType) *types.OmmType {
	var fn = ins.Fetch(funcname).Value
	var cargs = types.OmmArray{
		Array: args,
		Length: uint64(len(args)),
	} 

	return interpreter.Operations["function <- array"](*fn, cargs, ins, []string{ "at goat binder" }, 0, "goat")
}
