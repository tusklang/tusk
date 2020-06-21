package run

import "encoding/gob"
import "os"
import "fmt"

import "lang" //compiler
import . "lang/interpreter" //interpreter


//export Run
func Run(params map[string]map[string]interface{}) {
  dir := params["Files"]["DIR"]
  fileName := params["Files"]["NAME"]

  readfile, e := os.Open(dir.(string) + fileName.(string))

  if e != nil {
    fmt.Println("Error, could not access given oat file")
    os.Exit(1)
  }

  var decoded []Action

  decoder := gob.NewDecoder(readfile)
  e = decoder.Decode(&decoded)

  if e != nil {
    fmt.Println("Error, the given file is not oat compatible")
    os.Exit(1)
  }

  readfile.Close()

  //run the oat
  lang.OatRun(decoded, params, dir.(string))
}
