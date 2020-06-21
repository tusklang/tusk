package compile

import "encoding/gob"
import "os"

import "lang" //compiler

//export Compile
func Compile(params map[string]map[string]interface{}) {

  dir := params["FIles"]["DIR"]
  fileName := params["Files"]["NAME"]

  file := lang.ReadFileJS(dir.(string) + fileName.(string))[0]["Content"]

  lex := lang.Lexer(file, dir.(string), fileName.(string))
  acts := lang.Actionizer(lex, false, dir.(string), fileName.(string))

  if (IsAbsolute(params["Calc"]["O"].(string))) {

    writefile, _ := os.Create(params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(acts)
  } else {

    writefile, _ := os.Create(dir.(string) + params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(acts)
  }
}
