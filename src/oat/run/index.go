package oatRun

import "encoding/gob"
import "os"
import "fmt"

import ."lang/interpreter"

//export Run
func Run(params map[string]map[string]interface{}) {
  dir := params["Files"]["DIR"]
  fileName := params["Files"]["NAME"]

  readfile, e := os.Open(dir.(string) + fileName.(string))

  if e != nil {
    fmt.Println("Error, could not access given oat file")
    os.Exit(1)
  }

  var decoded map[string][]Action

  decoder := gob.NewDecoder(readfile)
  e = decoder.Decode(&decoded)

  if e != nil {
    fmt.Println("Error, the given file is not oat compatible")
    os.Exit(1)
  }

  readfile.Close()

  //run the oat
  RunInterpreter(decoded, params, dir.(string))
}
