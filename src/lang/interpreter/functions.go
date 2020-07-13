package interpreter

import "strconv"
import . "lang/types"

func callAsync(actions []Action, cli_params CliParams, stacktrace []string, ret chan Returner) {
  ret <- interpreter(actions, cli_params, stacktrace)
}

func initfuncs() {
  var function__sync__array = func(val1, val2 OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {
    var fn = val1.(OmmFunc)
    var arr = val2.(OmmArray)

    if uint64(len(fn.Params)) != arr.Length {
      ommPanic("Expected " + strconv.Itoa(len(fn.Params)) + " arguments to call function, but instead got " + strconv.FormatUint(arr.Length, 10), line, file, stacktrace)
    }

    for k, v := range arr.Array {
      vars[fn.Params[k]] = Variable{
        Type: "arg",
        Value: v,
      }
    }

    returnVal := interpreter(fn.Body, cli_params, append(stacktrace, "synchronous call at line " + strconv.FormatUint(line, 10) + " in file " + file))

    return returnVal.Exp
  }

  var function__async__array = func(val1, val2 OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {
    var fn = val1.(OmmFunc)
    var arr = val2.(OmmArray)

    if uint64(len(fn.Params)) != arr.Length {
      ommPanic("Expected " + strconv.Itoa(len(fn.Params)) + " arguments to call function, but instead got " + strconv.FormatUint(arr.Length, 10), line, file, stacktrace)
    }

    for k, v := range arr.Array {
      vars[fn.Params[k]] = Variable{
        Type: "arg",
        Value: v,
      }
    }

    channel := make(chan Returner)

    var promise OmmType = OmmThread{
      Channel: channel,
      Alive: true,
    }

    go callAsync(fn.Body, cli_params, append(stacktrace, "asynchronous call at line " + strconv.FormatUint(line, 10) + " in file " + file), channel)

    return &promise
  }

  operations["function sync array"], operations["function async array"] = function__sync__array, function__async__array
}
