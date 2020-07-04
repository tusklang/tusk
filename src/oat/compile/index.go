package oatCompile

import "encoding/gob"
import "os"

import "lang/compiler" //compiler
import . "lang/interpreter"

//export Compile
func Compile(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"]
  fileName := params["Files"]["NAME"]

  file := compiler.ReadFileJS(dir.(string) + fileName.(string))[0]["Content"]

  lex := compiler.Lexer(file, dir.(string), fileName.(string))
  _, variables := compiler.Actionizer(lex, false, dir.(string), fileName.(string))

  gob.Register(map[string][]Action{})

  if (IsAbsolute(params["Calc"]["O"].(string))) {

    writefile, _ := os.Create(params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(variables)
  } else {

    writefile, _ := os.Create(dir.(string) + params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(variables)
  }
}
