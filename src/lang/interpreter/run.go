package interpreter

import "os"
import "fmt"

import . "lang/types"

var threads []OmmThread
var vars map[string]Variable

//export RunInterpreter
func RunInterpreter(compiledVars map[string][]Action, cli_params map[string]map[string]interface{}) {

  vars = make(map[string]Variable)

  var dirnameOmmStr OmmString
  dirnameOmmStr.FromGoType(cli_params["Files"]["DIR"].(string))
  var dirnameOmmType OmmType = dirnameOmmStr

  vars["$__dirname"] = Variable{
    Type: "variable",
    Value: &dirnameOmmType,
  }

  initfuncs()

  for k, v := range compiledVars {
    vars[k] = Variable{
      Type: "variable",
      Value: interpreter(v, cli_params, []string{"at the global interpreter"}).Exp,
    }
  }

  for k, v := range GoFuncs {
    var gofunc OmmType = OmmGoFunc{
      Function: v,
    }
    vars["$" + k] = Variable{
      Type: "variable",
      Value: &gofunc,
    }
  }

  if _, exists := vars["$main"]; !exists {
    fmt.Println("Given program has no entry point/main function")
    os.Exit(1)
  } else {

    switch (*vars["$main"].Value).(type) {
      case OmmFunc:
        main := vars["$main"]

        calledP := interpreter((*main.Value).(OmmFunc).Body, cli_params, []string{"at the entry caller"}).Exp

        if calledP == nil {
          os.Exit(0)
        }

        called := *calledP

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
