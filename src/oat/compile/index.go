package compile

import "encoding/gob"
import "os"

import "lang" //omm language

//export Compile
func Compile(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"].(string)
  fileName := params["Files"]["NAME"].(string)

  file := lang.ReadFileJS(dir + fileName)[0]["Content"]

  lex := lang.Lexer(file, dir, fileName)
  acts := lang.Actionizer(lex, false, dir, fileName)

  if (IsAbsolute(params["Calc"]["O"].(string))) {

    writefile, _ := os.Create(params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(acts)
  } else {

    writefile, _ := os.Create(dir + params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(acts)
  }
}
