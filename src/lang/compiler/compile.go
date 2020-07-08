package compiler

import "encoding/json"
import "fmt"
import "path"

import . "lang/types"

//export Compile
func Compile(file, dir, filename string) []Action {

  lex := lexer(file, dir, filename)

  groups := makeGroups(lex)
  operations := makeOperations(groups)
  actions := actionizer(operations, path.Join(dir, filename))

  j, _ := json.MarshalIndent(actions, "", "  ")
  fmt.Println(string(j))

  return actions
}
