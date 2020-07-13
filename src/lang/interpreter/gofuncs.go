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
}
