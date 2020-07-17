package interpreter

import "strconv"
import . "lang/types"

func callAsync(actions []Action, instance *Instance, stacktrace []string, ret chan Returner) {
  ret <- instance.interpreter(actions, stacktrace)
}

func initfuncs() {
  var function__sync__array = func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {
    var fn = val1.(OmmFunc)
    var arr = val2.(OmmArray)

    if uint64(len(fn.Params)) != arr.Length {
      OmmPanic("Expected " + strconv.Itoa(len(fn.Params)) + " arguments to call function, but instead got " + strconv.FormatUint(arr.Length, 10), line, file, stacktrace)
    }

    for k, v := range arr.Array {
      instance.vars[fn.Params[k]] = &OmmVar{
        Name: fn.Params[k],
        Value: v,
      }
    }

    returnVal := instance.interpreter(fn.Body, append(stacktrace, "synchronous call at line " + strconv.FormatUint(line, 10) + " in file " + file))

    return returnVal.Exp
  }

  var function__async__array = func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {
    var fn = val1.(OmmFunc)
    var arr = val2.(OmmArray)

    if uint64(len(fn.Params)) != arr.Length {
      OmmPanic("Expected " + strconv.Itoa(len(fn.Params)) + " arguments to call function, but instead got " + strconv.FormatUint(arr.Length, 10), line, file, stacktrace)
    }

    for k, v := range arr.Array {
      instance.vars[fn.Params[k]] = &OmmVar{
        Name: fn.Params[k],
        Value: v,
      }
    }

    channel := make(chan Returner)

    var promise OmmType = OmmThread{
      Channel: channel,
    }

    go callAsync(fn.Body, instance, append(stacktrace, "asynchronous call at line " + strconv.FormatUint(line, 10) + " in file " + file), channel)

    return &promise
  }

  var gofunc__sync__array = func(val1, val2 OmmType, instance *Instance, stacktrace []string, line uint64, file string) *OmmType {
    gfn := val1.(OmmGoFunc)
    arr := val2.(OmmArray)

    return gfn.Function(arr.Array, stacktrace, line, file)
  }

  operations["function <- array"] = function__sync__array
  operations["function <~ array"] = function__async__array
  operations["gofunc <- array"] = gofunc__sync__array
}
