package oat

import "os"
import "fmt"
import "io/ioutil"
// import "encoding/gob"
// import "bytes"
import "strings"

import . "lang/types"
import . "oat/encoding"
import "lang/compiler" //compiler

//export Compile
func Compile(params CliParams) {
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
  _, vars, ce := compiler.Compile(string(file), fileName, compileall, true)
  compiler.Ommbasedir = "" //reset Ommbasedir

  if ce != nil {
    ce.Print()
  }

  // var vals = Oat{
  //   Actions: actions,
  //   Variables: vars,
  // }

  OatEncode(params.Output, vars)

  // var network bytes.Buffer
  // encoder := gob.NewEncoder(&network)

  // if err := encoder.Encode(vals); err != nil {
  //   panic(err)
  // }

  // nbytes := network.Bytes()
  // writer, e := os.Create(params.Output)

  // if e != nil {
  //   fmt.Println("Could not make file:", fileName)
  //   os.Exit(1)
  // }

  // _, e = writer.Write(nbytes)

  // if e != nil {
  //   fmt.Println("Could not write file:", fileName)
  //   os.Exit(1)
  // }
}
