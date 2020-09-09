package interpreter

import (
	"strconv"

	. "github.com/omm-lang/omm/lang/types"
	"github.com/omm-lang/omm/native"
)

//fill an instance to run a function
func fillFuncInstance(fn *OmmFunc, args OmmArray, parent *Instance) *Overload {

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
	var function__sync__array = func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		var fn = val1.(OmmFunc)
		var arr = val2.(OmmArray)

		var overload *Overload

		if overload = fillFuncInstance(&fn, arr, instance); overload == nil {
			native.OmmPanic("Could not find a typelist for function call", line, file, stacktrace)
		}

		return Interpreter(fn.Instance, overload.Body, append(stacktrace, "asynchronous call at line "+strconv.FormatUint(line, 10)+" in file "+file), stacksize+1, overload.Params, true).Exp
	}

	var function__async__array = func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		var fn = val1.(OmmFunc)
		var arr = val2.(OmmArray)

		var overload *Overload

		if overload = fillFuncInstance(&fn, arr, instance); overload == nil {
			native.OmmPanic("Could not find a typelist for function call", line, file, stacktrace)
		}

		var promise OmmType = *NewThread(func() *OmmType {
			return Interpreter(fn.Instance, overload.Body, append(stacktrace, "asynchronous call at line "+strconv.FormatUint(line, 10)+" in file "+file), stacksize+1, overload.Params, true).Exp
		})

		return &promise
	}

	var nativefunc__sync__array = func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		gfn := val1.(native.OmmGoFunc)
		arr := val2.(OmmArray)

		if gfn.Function == nil {
			native.OmmPanic("Native function is nil", line, file, stacktrace)
		}

		return gfn.Function(arr.Array, stacktrace, line, file, instance)
	}

	Operations["function : array"] = function__sync__array
	Operations["function ? array"] = function__async__array
	Operations["native_func : array"] = nativefunc__sync__array
}
