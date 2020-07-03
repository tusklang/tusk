package oatCompile

import "encoding/gob"
import "os"

import "lang/compiler" //compiler
import . "lang/interpreter"

type OatEncode struct {
  Actions []Action
  Compiled_Variables map[string][]Action
}

//export Compile
func Compile(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"]
  fileName := params["Files"]["NAME"]

  file := compiler.ReadFileJS(dir.(string) + fileName.(string))[0]["Content"]

  lex := compiler.Lexer(file, dir.(string), fileName.(string))
  acts, variables := compiler.Actionizer(lex, false, dir.(string), fileName.(string))

  gob.Register(map[string]interface{}{})

  var putValue = OatEncode{
    Actions: acts,
    Compiled_Variables: variables,
  }

  if (IsAbsolute(params["Calc"]["O"].(string))) {

    writefile, _ := os.Create(params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(putValue)
  } else {

    writefile, _ := os.Create(dir.(string) + params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(putValue)
  }
}
