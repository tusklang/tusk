package oat

import "os"
import "fmt"
import "io/ioutil"
import "encoding/gob"
import "bytes"
import "oat/helper"
import "strings"

import . "lang/types"
import "lang/compiler" //compiler

//export Compile
func Compile(params CliParams) {
  oatHelper.InitGob()

  fileName := params.Name

  file, e := ioutil.ReadFile(fileName)

  if e != nil {
    fmt.Println("Could not find file:", fileName)
    os.Exit(1)
  }

  var compileall = false
  if strings.HasSuffix(fileName, "*") || strings.HasSuffix(fileName, "*/") {
    compileall = true
    fileName = "main.omm"
  }

  compiler.Ommbasedir = params.OmmDirname
  actions, vars, ce := compiler.Compile(string(file), fileName, compileall, true)
  compiler.Ommbasedir = "" //reset Ommbasedir

  if ce != nil {
    ce.Print()
  }

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
  writer, e := os.Create(params.Output)

  if e != nil {
    fmt.Println("Could not make file:", fileName)
    os.Exit(1)
  }

  _, e = writer.Write(nbytes)

  if e != nil {
    fmt.Println("Could not write file:", fileName)
    os.Exit(1)
  }
}
