package compiler

import "encoding/json"
import "fmt"
import "path"
import "os"
import "strconv"

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

  j, _ := json.MarshalIndent(actions, "", "  ")
  fmt.Println(string(j))

  return actions, make(map[string][]Action)
}
