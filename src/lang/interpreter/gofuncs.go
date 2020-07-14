package interpreter

//all of the gofuncs
//functions written in go that are used by omm

import "bufio"
import "os"
import "fmt"
import "time"

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
          OmmPanic("Expected a string as the argument to input[]", line, file, stacktrace)
      }

    } else {
      OmmPanic("Function input requires a parameter count of 0 or 1", line, file, stacktrace)
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
      OmmPanic("Function typeof requires a parameter count of 1", line, file, stacktrace)
    }

    typeof := (*args[0]).TypeOf()

    var str OmmString
    str.FromGoType(typeof)

    //convert to OmmType interface
    var ommtype OmmType = str

    return &ommtype
  },
  "defop": func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {

    if len(args) != 4 {
      OmmPanic("Function defop requires a parameter count of 4", line, file, stacktrace)
    }

    if (*args[0]).Type() != "string" || (*args[1]).Type() != "string" || (*args[2]).Type() != "string" || (*args[3]).Type() != "function" {
      OmmPanic("Function defop requires [string, string, string, function]", line, file, stacktrace)
    }

    operation := (*args[0]).(OmmString).ToGoType()
    operand1 := (*args[1]).(OmmString).ToGoType()
    operand2 := (*args[2]).(OmmString).ToGoType()
    function := (*args[3]).(OmmFunc)

    if len(function.Params) != 2 {
      OmmPanic("Expected a parameter count of 2 for the fourth argument of defop", line, file, stacktrace)
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
      OmmPanic("Function append requires a parameter count of 2", line, file, stacktrace)
    }

    if (*args[0]).Type() != "array" {
      OmmPanic("Function append requires the first argument to be an array", line, file, stacktrace)
    }

    appended := append((*args[0]).(OmmArray).Array, args[1])
    var arr OmmType = OmmArray{
      Array: appended,
      Length: uint64(len(appended)),
    }

    return &arr
  },
  "exit": func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {

    if len(args) == 1 {

      switch (*args[0]).(type) {
        case OmmNumber:

          var gonum = (*args[0]).(OmmNumber).ToGoType()
          os.Exit(int(gonum))

        case OmmBool:

          if (*args[0]).(OmmBool).ToGoType() == true {
            os.Exit(0)
          } else {
            os.Exit(1)
          }

        default:
          os.Exit(0)
      }

    } else if len(args) == 0 {
      os.Exit(0)
    } else {
      OmmPanic("Function exit requires a parameter count of 1 or 0", line, file, stacktrace)
    }

    var tmpundef OmmType = undef
    return &tmpundef
  },
  "wait": func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {

    if len(args) == 1 {

      if (*args[0]).Type() != "number" {
        OmmPanic("Function wait requires a number as the argument", line, file, stacktrace)
      }

      var amt = (*args[0]).(OmmNumber)

      var n4294967295 = zero
      n4294967295.Integer = &[]int64{ 5, 9, 2, 7, 6, 9, 4, 9, 2, 4 }

      //if amt is less than 2 ^ 32 - 1, just convert to a go int
      if isLess(amt, n4294967295) {
        gonum := amt.ToGoType()

        time.Sleep(time.Duration(gonum) * time.Millisecond)
      } else {
        //this is how this works
        /*
          the loop starts at 0 with an increment of 2 ^ 32 - 1
          in each iteration, it will wait for 2 ^ 32 - 1 milliseconds
        */

        for i := zero; isLess(i, amt); i = (*number__plus__number(i, n4294967295, cli_params, stacktrace, line, file)).(OmmNumber) {
          time.Sleep(4294967295 * time.Millisecond)
        }
      }

    } else {
      OmmPanic("Function wait requires a parameter count of 1", line, file, stacktrace)
    }

    var tmpundef OmmType = undef
    return &tmpundef
  },
  "thread.wasjoined": func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {

    if len(args) == 1 {

      switch (*args[0]).(type) {
        case OmmThread:

          var wasJoined = (*args[0]).(OmmThread).WasJoined
          var ommtype OmmType = OmmBool{
            Boolean: &wasJoined,
          }

          return &ommtype
        default:
          OmmPanic("Function thread.wasjoined requires a thread as the argument", line, file, stacktrace)
      }

    } else {
      OmmPanic("Function thread.wasjoined requires a parameter count of 1", line, file, stacktrace)
    }

    var tmpfalse OmmType = falsev
    return &tmpfalse
  },
  "make": func(args []*OmmType, cli_params CliParams, stacktrace []string, line uint64, file string) *OmmType {

    if len(args) == 1 {

      switch (*args[0]).(type) {
      case OmmProto:
          var ommtype OmmType = OmmObject{}.New((*args[0]).(OmmProto))
          return &ommtype
        default:
          OmmPanic("Function make requires a structure as the argument", line, file, stacktrace)
      }

    } else {
      OmmPanic("Function make requires a parameter count of 1", line, file, stacktrace)
    }

    var tmpundef OmmType = undef
    return &tmpundef
  },
}
