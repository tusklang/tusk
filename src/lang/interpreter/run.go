package interpreter

import "math"

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
func RunInterpreter(actions []Action, cli_params map[string]map[string]interface{}, dir string) {

  var vars = make(map[string]Variable)

  for k, v := range goprocs {
    vars["$" + k] = Variable{
      Type: "goproc",
      Name: "$" + k,
      GoProc: v,
    }
  }

  interpreter(actions, CliParams(cli_params), vars, false, []Action{}, dir)

  for _, v := range threads {
    v.WaitFor()
  }
}
