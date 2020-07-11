package interpreter

import "os"
import "fmt"

import . "lang/types"

var threads []OmmThread

//export RunInterpreter
func RunInterpreter(compiledVars map[string][]Action, cli_params map[string]map[string]interface{}) {

  var vars = make(map[string]Variable)

  for k, v := range compiledVars {
    vars[k] = Variable{
      Type: "global",
      Value: interpreter(v, cli_params, vars, make([]Action, 0)).Exp,
    }
  }

  for k, v := range GoFuncs {
    vars["$" + k] = Variable{
      Type: "gofunc",
      GoProc: v,
    }
  }

  if _, exists := vars["$main"]; !exists {
    fmt.Println("Given program has no entry point/main function")
    os.Exit(1)
  } else {

    switch vars["$main"].Value.(type) {
      case OmmFunc:
        main := vars["$main"]

        called := interpreter(main.Value.(OmmFunc).Body, cli_params, vars, make([]Action, 0)).Exp

        for _, v := range threads {
          v.WaitFor()
        }

        var exitType int64 = 0

        switch called.(type) {
          case OmmNumber:
            exitType = int64(called.(OmmNumber).ToGoType())
        }

        os.Exit(int(exitType))
      default:
        fmt.Println("Entry point was not given as a function")
        os.Exit(1)
    }
  }
}
