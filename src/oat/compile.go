package oat

import "os"
import "fmt"
import "io/ioutil"
import "encoding/gob"
import "bytes"
import "oat/helper"

import . "lang/types"
import "lang/compiler" //compiler

//export Compile
func Compile(params CliParams) {
  oatHelper.InitGob()

  dir := params["Files"]["DIR"]
  fileName := params["Files"]["NAME"]

  file, e := ioutil.ReadFile(fileName.(string))

  if e != nil {
    fmt.Println("Could not find file:", fileName.(string))
    os.Exit(1)
  }

  actions, vars := compiler.Compile(string(file), dir.(string), fileName.(string))

  var vals = Oat{
    Actions: actions,
    Variables: vars,
  }

  var network bytes.Buffer
  encoder := gob.NewEncoder(&network)

  if err := encoder.Encode(vals); err != nil {
    panic(err)
  }

  nbytes := network.Bytes()
  writer, e := os.Create(params["Calc"]["O"].(string))

  if e != nil {
    fmt.Println("Could not make file:", fileName.(string))
    os.Exit(1)
  }

  _, e = writer.Write(nbytes)

  if e != nil {
    fmt.Println("Could not write file:", fileName.(string))
    os.Exit(1)
  }
}
