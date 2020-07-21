package interpreter

import "os"
import "fmt"

import . "lang/types"

var threads []OmmThread

//export RunInterpreter
func RunInterpreter(compiledVars map[string][]Action, cli_params CliParams) {

  var dirnameOmmStr OmmString
  dirnameOmmStr.FromGoType(cli_params.Directory)
  var dirnameOmmType OmmType = dirnameOmmStr

  var instance Instance

  instance.Params = cli_params
  instance.vars = make(map[string]*OmmVar)
  instance.globals = make(map[string]*OmmVar)

  instance.vars["$__dirname"] = &OmmVar{
    Name: "$__dirname",
    Value: &dirnameOmmType,
  }

  for k, v := range compiledVars {
    var global = OmmVar{
      Name: k,
      Value: instance.interpreter(v, []string{"at the global interpreter"}).Exp,
    }
    instance.globals[k] = &global
    instance.vars = instance.globals //copy globals every time because some globals reference each other
  }

  for k, v := range GoFuncs {
    var gofunc OmmType = OmmGoFunc{
      Function: v,
    }
    instance.vars["$" + k] = &OmmVar{
      Name: "$" + k,
      Value: &gofunc,
    }
  }

  if _, exists := instance.vars["$main"]; !exists {
    fmt.Println("Given program has no entry point/main function")
    os.Exit(1)
  } else {

    switch (*instance.vars["$main"].Value).(type) {
      case OmmFunc:
        main := instance.vars["$main"]

        calledP := instance.interpreter((*main.Value).(OmmFunc).Body, []string{"at the entry caller"}).Exp

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
