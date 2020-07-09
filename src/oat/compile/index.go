package oatCompile

import "encoding/gob"
import "os"

import "lang/compiler" //compiler

//export Compile
func Compile(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"]
  fileName := params["Files"]["NAME"]

  file := compiler.ReadFileJS(dir.(string) + fileName.(string))[0]["Content"]

  actions, vars := compiler.Compile(file, dir.(string), fileName.(string))

  var vals = compiler.OatValues{ actions, vars, params }

  if (IsAbsolute(params["Calc"]["O"].(string))) {

    writefile, _ := os.Create(params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(vals)
  } else {

    writefile, _ := os.Create(dir.(string) + params["Calc"]["O"].(string))

    defer writefile.Close()

    encoder := gob.NewEncoder(writefile)
    encoder.Encode(vals)
  }
}
