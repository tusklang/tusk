package compiler

import "encoding/json"
import "fmt"

import . "lang/interpreter"

//export Compile
func Compile(file, dir, filename string) []Action {

  lex := lexer(file, dir, filename)
  groups := makeGroups(lex)

  j, _ := json.MarshalIndent(groups, "", "  ")
  fmt.Println(string(j))

  return []Action{}
}
