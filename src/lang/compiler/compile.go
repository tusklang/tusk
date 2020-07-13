package compiler

import "encoding/json"
import "fmt"
import "path"
import "os"
import "strconv"

import "lang/interpreter"
import . "lang/types"

type OatValues struct {
  Actions    []Action
  Variables    map[string][]Action
  Params       map[string]map[string]interface{}
}

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
func Compile(file, dir, filename string) ([]Action, map[string][]Action) {

  lex := lexer(file, dir, filename)
  groups := makeGroups(lex)
  operations := makeOperations(groups)
  actions := actionizer(operations, path.Join(dir, filename))
  vars := getvars(actions, path.Join(dir, filename))

  //make each var have only it's name
  var varnames = make(map[string]string)

  for k := range vars {
    varnames[k] = k
  }

  //also account for the gofuncs
  for k := range interpreter.GoFuncs {
    varnames["$" + k] = "$" + k
  }

  for k := range vars {
    changevarnames(vars[k], varnames)
  }

  j, _ := json.MarshalIndent(vars, "", "  ")
  fmt.Println(string(j))

  return actions, vars
}
