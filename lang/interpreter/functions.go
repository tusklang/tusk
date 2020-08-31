package interpreter

import (
	"strconv"

	. "github.com/omm-lang/omm/lang/types"
	"github.com/omm-lang/omm/ommstd/native"
	. "github.com/omm-lang/omm/ommstd/native"
)

func funcinit() { //initialize the operations that require the use of the interpreter
	var function__sync__array = func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		var fn = val1.(OmmFunc)
		var arr = val2.(OmmArray)

		for _, v := range fn.Overloads {
			if uint64(len(v.Params)) != arr.Length {
				continue
			}

			var not_exists bool = false

			for k, vv := range v.Types {
				if vv != (*arr.At(int64(k))).TypeOf() && vv != "any" {
					not_exists = true
					break
				}
			}

			if !not_exists {

				if fn.Instance == nil {
					fn.Instance = instance
				}

				for k, vv := range arr.Array {
					fn.Instance.Allocate(v.Params[k], vv)
				}

				return Interpreter(fn.Instance, v.Body, append(stacktrace, "synchronous call at line "+strconv.FormatUint(line, 10)+" in file "+file), stacksize+1, v.Params).Exp
			}
		}

		OmmPanic("Could not find a typelist for function call", line, file, stacktrace)

		var tmp OmmType = undef
		return &tmp
	}

	var function__async__array = func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		var fn = val1.(OmmFunc)
		var arr = val2.(OmmArray)

		for _, v := range fn.Overloads {
			if uint64(len(v.Params)) != arr.Length {
				continue
			}

			var not_exists bool = false

			for k, vv := range v.Types {
				if vv != (*arr.At(int64(k))).TypeOf() && vv != "any" {
					not_exists = true
					break
				}
			}

			if !not_exists {

				if fn.Instance == nil {
					fn.Instance = instance
				}

				for k, vv := range arr.Array {
					fn.Instance.Allocate(v.Params[k], vv)
				}

				var promise OmmType = *NewThread(func() *OmmType {
					return Interpreter(fn.Instance, v.Body, append(stacktrace, "asynchronous call at line "+strconv.FormatUint(line, 10)+" in file "+file), stacksize+1, v.Params).Exp
				})

				return &promise
			}
		}

		OmmPanic("Could not find a typelist for function call", line, file, stacktrace)

		var tmp OmmType = undef
		return &tmp
	}

	var gofunc__sync__array = func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string, stacksize uint) *OmmType {
		gfn := val1.(OmmGoFunc)
		arr := val2.(OmmArray)

		if gfn.Function == nil {
			native.OmmPanic("Native function is nil", line, file, stacktrace)
		}

		return gfn.Function(arr.Array, stacktrace, line, file, instance)
	}

	Operations["function <- array"] = function__sync__array
	Operations["function <~ array"] = function__async__array
	Operations["gofunc <- array"] = gofunc__sync__array
}
