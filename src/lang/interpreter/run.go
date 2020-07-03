package interpreter

import "os"
import "fmt"
import "math"
import "strconv"

type CliParams map[string]map[string]interface{}

//number sizes

//export DigitSize
const DigitSize = 1;
//export MAX_DIGIT
var MAX_DIGIT = int64(math.Pow(10, DigitSize) - 1)
//export MIN_DIGIT
var MIN_DIGIT = -1 * MAX_DIGIT

//////////////

//export RunInterpreter
func RunInterpreter(actions []Action, cli_params map[string]map[string]interface{}, dir string, compiledVars map[string][]Action) {
  var vars = make(map[string]Variable)

  for k, v := range compiledVars {
    vars[k] = Variable{
      Type: "global",
      Name: k,
      Value: interpreter(v, cli_params, vars, true, make([]Action, 0), dir).Exp,
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
    called := interpreter(main.Value.ExpAct, cli_params, vars, true, make([]Action, 0), dir).Exp
    exitType, _ := strconv.Atoi(cast(called, "string").ExpStr) //convert return value to int
    os.Exit(exitType)
  }

  for _, v := range threads {
    v.WaitFor()
  }
}
