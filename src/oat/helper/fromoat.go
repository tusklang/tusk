package oatHelper

import "os"
import "encoding/gob"
import "fmt"
import . "lang/types"

//export FromOat
func FromOat(fileName string) Oat {
  InitGob()

  readfile, e := os.Open(fileName)

  if e != nil {
    fmt.Println("Error, could not access given oat file:", fileName)
    os.Exit(1)
  }

  var decoded Oat

  decoder := gob.NewDecoder(readfile)
  e = decoder.Decode(&decoded)

  if e != nil {
    fmt.Println("Error, the given file is not oat compatible:", fileName)
    os.Exit(1)
  }

  readfile.Close()
  return decoded
}
