package interpreter

import (
	"fmt"
	"strconv"

	. "github.com/tusklang/tusk/lang/types"
	. "github.com/tusklang/tusk/native"
)

//fill an instance to run a function
func fillFuncInstance(fn *TuskFunc, args TuskArray, parent *Instance) *Overload {

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
	var function__sync__array = func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		var fn = val1.(TuskFunc)
		var arr = val2.(TuskArray)

		var overload *Overload

		if overload = fillFuncInstance(&fn, arr, instance); overload == nil {
			return nil, TuskPanic("Could not find a typelist for function call", line, file, stacktrace)
		}

		tmp, e := Interpreter(fn.Instance, overload.Body, append(stacktrace, "synchronous call at line "+strconv.FormatUint(line, 10)+" in file "+file), stacksize+1, overload.Params, true)
		return tmp.Exp, e
	}

	var function__async__array = func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		var fn = val1.(TuskFunc)
		var arr = val2.(TuskArray)

		var overload *Overload

		if overload = fillFuncInstance(&fn, arr, instance); overload == nil {
			return nil, TuskPanic("Could not find a typelist for function call", line, file, stacktrace)
		}

		var promise TuskType = *NewThread(func() (*TuskType, *TuskError) {
			tmp, e := Interpreter(fn.Instance, overload.Body, append(stacktrace, "asynchronous call at line "+strconv.FormatUint(line, 10)+" in file "+file), stacksize+1, overload.Params, true)
			return tmp.Exp, e
		})

		return &promise, nil
	}

	var nativefunc__sync__array = func(val1, val2 TuskType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) (*TuskType, *TuskError) {
		gfn := val1.(TuskGoFunc)
		arr := val2.(TuskArray)

		if gfn.Function == nil {
			return nil, TuskPanic("Native function is nil", line, file, stacktrace)
		}

		//check the signatures of the function

		var sigmatch = true

		for _, v := range gfn.Signatures {
			//if it is {"..."} it works no matter what
			if len(v) != 0 && v[0] == "..." {
				break
			}

			{
				for kk, vv := range v {
					curv := *arr.At(int64(kk))
					if vv == "any" || curv.TypeOf() == vv || curv.Type() == vv {
						//all of these are good
						sigmatch = true
						continue
					} else {
						//otherwise wrong sig
						sigmatch = false
						break
					}
				}
			}

		}

		if sigmatch { //if a signature matches
			return gfn.Function(arr.Array, stacktrace, line, file, instance)
		}

		//otherwise panic
		return nil, TuskPanic(fmt.Sprintf("Native function requires the signature %s", func() string {

			//give the signature list

			var ret string
			for k, v := range gfn.Signatures {
				ret += "("
				for kk, vv := range v {
					ret += vv
					if kk+1 != len(v) {
						ret += ", "
					}
				}
				ret += ")"

				if k+1 != len(gfn.Signatures) {
					ret += " or "
				}
			}
			return ret
		}()), line, file, stacktrace)
		//////////////////////////////////////

	}

	Operations["function : array"] = function__sync__array
	Operations["function ? array"] = function__async__array
	Operations["native_func : array"] = nativefunc__sync__array
}
