package compiler

import "path"
import "io/ioutil"
import "fmt"
import "os"

import . "lang/interpreter"

var included = []string{} //list of the imported files from omm

//export Run
func Run(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"]
  fileName := params["Files"]["NAME"]

  included = append(included, path.Join(dir.(string), fileName.(string)))

  file, e := ioutil.ReadFile(fileName.(string))

  if e != nil {
    fmt.Println("Could not find", fileName.(string))
    os.Exit(1)
  }

  _, variables := Compile(string(file), dir.(string), fileName.(string))

  RunInterpreter(variables, params)
}
