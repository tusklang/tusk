package interpreter

import "os"
import "fmt"

import . "lang/types"

var threads []OmmThread

//export RunInterpreter
func RunInterpreter(compiledVars map[string][]Action, cli_params map[string]map[string]interface{}) {
  var vars = make(map[string]Variable)

  for k, v := range compiledVars {
    _ = v
    vars[k] = Variable{
      Type: "global",
      Name: k,
      // Value: interpreter(v, cli_params, vars, true, make([]Action, 0), dir).Exp,
    }
  }

  for k, v := range gofuncs {
    vars["$" + k] = Variable{
      Type: "gofunc",
      Name: "$" + k,
      GoProc: v,
    }
  }

  if _, exists := vars["$main"]; !exists || vars["$main"].Value.Type != "function" {
    fmt.Println("Given program has no entry point/main function")
    os.Exit(1)
  } else {
    main := vars["$main"]
    _ = main
    // called := interpreter(main.Value.ExpAct, cli_params, vars, true, make([]Action, 0), dir).Exp

    for _, v := range threads {
      v.WaitFor()
    }

    // exitType, _ := strconv.Atoi(cast(called, "string").ExpStr) //convert return value to int
    // os.Exit(exitType)
  }
}
