package interpreter

import (
	"strconv"

	. "ka/lang/types"
)

//fill an instance to run a function
func fillFuncInstance(fn *KaFunc, args KaArray, parent *Instance) *Overload {

	if fn.Instance == nil {
		fn.Instance = (*parent).Copy() //copy the parent instance, if it is not part of an object
	}

	for _, v := range fn.Overloads {
		if uint64(len(v.Params)) != args.Length {
			continue
		}

		for k, vv := range v.Types {
			if vv != (*args.At(int64(k))).TypeOf() && vv != "any" {
				goto wrong_signature
				break
			}
		}

		{
			for k, vv := range args.Array {
				fn.Instance.Allocate(v.Params[k], vv)
			}
			return &v
		}

	wrong_signature: //the current signature doesn't match
	}

	return nil
}

func funcinit() { //initialize the operations that require the use of the interpreter
	var function__sync__array = func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		var fn = val1.(KaFunc)
		var arr = val2.(KaArray)

		var overload *Overload

		if overload = fillFuncInstance(&fn, arr, instance); overload == nil {
			KaPanic("Could not find a typelist for function call", line, file, stacktrace)
		}

		return Interpreter(fn.Instance, overload.Body, append(stacktrace, "asynchronous call at line "+strconv.FormatUint(line, 10)+" in file "+file), stacksize+1, overload.Params, true).Exp
	}

	var function__async__array = func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		var fn = val1.(KaFunc)
		var arr = val2.(KaArray)

		var overload *Overload

		if overload = fillFuncInstance(&fn, arr, instance); overload == nil {
			KaPanic("Could not find a typelist for function call", line, file, stacktrace)
		}

		var promise KaType = *NewThread(func() *KaType {
			return Interpreter(fn.Instance, overload.Body, append(stacktrace, "asynchronous call at line "+strconv.FormatUint(line, 10)+" in file "+file), stacksize+1, overload.Params, true).Exp
		})

		return &promise
	}

	var nativefunc__sync__array = func(val1, val2 KaType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *KaType {
		gfn := val1.(KaGoFunc)
		arr := val2.(KaArray)

		if gfn.Function == nil {
			KaPanic("Native function is nil", line, file, stacktrace)
		}

		return gfn.Function(arr.Array, stacktrace, line, file, instance)
	}

	Operations["function : array"] = function__sync__array
	Operations["function ? array"] = function__async__array
	Operations["native_func : array"] = nativefunc__sync__array
}
