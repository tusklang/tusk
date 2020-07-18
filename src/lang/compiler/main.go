package compiler

import "io/ioutil"
import "fmt"
import "os"

import . "lang/types"
import . "lang/interpreter"

var included = []string{} //list of the imported files from omm

//export Run
func Run(params CliParams) {

  fileName := params.Name

  included = append(included, fileName)

  file, e := ioutil.ReadFile(fileName)

  if e != nil {
    fmt.Println("Could not find", fileName)
    os.Exit(1)
  }

  _, variables, ce := Compile(string(file), fileName)

  if ce != nil {
    ce.Print()
  }

  RunInterpreter(variables, params)
}
