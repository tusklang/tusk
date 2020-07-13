package interpreter

//all of the gofuncs
//functions written in go that are used by omm

import "bufio"
import "os"
import "fmt"

import . "lang/types"

//export GoFuncs
var GoFuncs = map[string]func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {
  "input": func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {

    scanner := bufio.NewScanner(os.Stdin)

    if len(args) == 0 {
      //if it has 0 or 1 arg, there is no error
    } else if len(args) == 1 {

      switch (*args[0]).(type) {
        case OmmString:
          str := (*args[0]).(OmmString).ToGoType()
          fmt.Print(str + ": ")
        default:
          ommPanic("Expected a string as the argument to input[]", line, file, stacktrace)
      }

    } else {
      ommPanic("Function input requires a parameter count of 0 or 1", line, file, stacktrace)
    }

    //get user input and convert it to OmmType
    scanner.Scan()
    input := scanner.Text()
    var inputOmmType OmmString
    inputOmmType.FromGoType(input)
    var inputType OmmType = inputOmmType

    return &inputType
  },
  "typeof": func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {

    if len(args) != 1 {
      ommPanic("Function typeof requires a parameter count of 1", line, file, stacktrace)
    }

    typeof := (*args[0]).Type()

    var str OmmString
    str.FromGoType(typeof)

    //convert to OmmType interface
    var ommtype OmmType = str

    return &ommtype
  },
  "defop": func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {

    if len(args) != 4 {
      ommPanic("Function defop requires a parameter count of 4", line, file, stacktrace)
    }

    if (*args[0]).Type() != "string" || (*args[1]).Type() != "string" || (*args[2]).Type() != "string" || (*args[3]).Type() != "function" {
      ommPanic("Function defop requires [string, string, string, function]", line, file, stacktrace)
    }

    operation := (*args[0]).(OmmString).ToGoType()
    operand1 := (*args[1]).(OmmString).ToGoType()
    operand2 := (*args[2]).(OmmString).ToGoType()
    function := (*args[3]).(OmmFunc)

    if len(function.Params) != 2 {
      ommPanic("Expected a parameter count of 2 for the fourth argument of defop", line, file, stacktrace)
    }

    operations[operand1 + " " + operation + " " + operand2] = func(val1, val2 OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {
      vars[function.Params[0]] = Variable{
        Type: "arg",
        Value: &val1,
      }
      vars[function.Params[1]] = Variable{
        Type: "arg",
        Value: &val2,
      }

      return interpreter(function.Body, cli_params, stacktrace).Exp
    }

    var tmpundef OmmType = undef
    return &tmpundef //return undefined
  },
  "append": func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {

    if len(args) != 2 {
      ommPanic("Function append requires a parameter count of 2", line, file, stacktrace)
    }

    if (*args[0]).Type() != "array" {
      ommPanic("Function append requires the first argument to be an array", line, file, stacktrace)
    }

    appended := append((*args[0]).(OmmArray).Array, args[1])
    var arr OmmType = OmmArray{
      Array: appended,
      Length: uint64(len(appended)),
    }

    return &arr
  },
}
