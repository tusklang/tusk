package compile

import "encoding/json"

import "lang" //omm language

// #cgo CFLAGS: -std=c99
// #include "bind.h"
import "C"

//export Compile
func Compile(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"].(string)
  fileName := params["Files"]["NAME"].(string)

  file := lang.ReadFileJS(dir + fileName)[0]

  lex := lang.Lexer(file, dir, fileName)
  acts := lang.Actionizer(lex, false, dir, fileName)

  var encoding string

  _encoding, _ := json.Marshal(acts)
  encoding = string(_encoding)

  //write oat file
  C.write(C.CString(dir), C.CString(params["Calc"]["O"].(string)), C.CString(encoding))
}
