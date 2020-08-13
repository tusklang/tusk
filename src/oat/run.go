package oat

import "fmt"
import "os"
import "lang/interpreter"
import . "oat/encoding"
import . "lang/types"

//export Run
func Run(params CliParams) {
  d, e := OatDecode(params.Name, 0)
  if e != nil {
    fmt.Println(e)
    os.Exit(1)
  }

  interpreter.RunInterpreter(d, params)
}
