package compiler

import "io/ioutil"
import "fmt"
import "os"

import . "lang/types"
import . "lang/interpreter"

var included = []string{} //list of the imported files from omm

//export Run
func Run(params CliParams) {

  fileName := params["Files"]["NAME"]

  included = append(included, fileName.(string))

  file, e := ioutil.ReadFile(fileName.(string))

  if e != nil {
    fmt.Println("Could not find", fileName.(string))
    os.Exit(1)
  }

  _, variables, ce := Compile(string(file), fileName.(string))

  if ce != nil {
    ce.Print()
  }

  RunInterpreter(variables, params)
}
