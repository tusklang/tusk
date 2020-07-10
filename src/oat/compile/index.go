package oatCompile

import "encoding/gob"
import "os"
import "fmt"
import "io/ioutil"

import "lang/compiler" //compiler

//export Compile
func Compile(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"]
  fileName := params["Files"]["NAME"]

  file, e := ioutil.ReadFile(fileName.(string))

  if e != nil {
    fmt.Println("Could not find file:", fileName.(string))
    os.Exit(1)
  }

  actions, vars := compiler.Compile(string(file), dir.(string), fileName.(string))

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
