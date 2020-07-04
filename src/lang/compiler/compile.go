package compiler

import "encoding/json"
import "fmt"

import . "lang/interpreter"

//export Compile
func Compile(file, dir, filename string) []Action {

  lex := lexer(file, dir, filename)

  j, _ := json.MarshalIndent(lex, "", "  ")
  fmt.Println(string(j))

  return []Action{}
}
