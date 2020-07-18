package compiler

import "fmt"
import "os"
import "strconv"

import "lang/interpreter"
import . "lang/types"

func compilerErr(msg string, dir string, line uint64) {

  //I dont know why regular printing doesnt work
  //(I am on ubuntu, but on windows it works)
  //(I just switched from windows to ubuntu)
  fmt.Print("Error while compiling " + dir)
  fmt.Println()
  fmt.Print("Error on line " + strconv.FormatUint(line, 10))
  fmt.Println("\n")
  fmt.Println(msg)
  //////////////////////////////////////////////

  os.Exit(1)
}

//export Compile
func Compile(file, filename string) ([]Action, map[string][]Action) {

  lex := lexer(file, filename)
  groups := makeGroups(lex)
  operations := makeOperations(groups)
  actions := actionizer(operations, filename)

  //a bunch of validations and initializers
  has_non_global_prototypes(actions, true)
  put_proto_types(actions)
  validate_types(actions)
  /////////////////////////////////////////

  vars := getvars(actions, filename)

  //make each var have only it's name
  var varnames = make(map[string]string)

  for k := range vars {

    if vars[k][0].Type != "function" {
      changevarnames(vars[k], varnames) //ensure none of the globals use the globals from below
    }
    varnames[k] = k
  }

  //also account for the gofuncs
  for k := range interpreter.GoFuncs {
    varnames["$" + k] = "$" + k
  }


  for k := range vars {
    changevarnames(vars[k], varnames)
  }

  return actions, vars
}
