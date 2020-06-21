package run

import "encoding/json"
import "encoding/gob"
import "os"
import "fmt"
import "lang/interpreter/bind"

import "lang" //omm language

//export Run
func Run(params bind.CliParams) {
  dir := params.GetFiles().GetDIR()
  fileName := params.GetFiles().GetNAME()

  readfile, e := os.Open(dir + fileName)

  if e != nil {
    fmt.Println("Error, could not access given oat file")
    os.Exit(1)
  }

  var decoded []lang.Action

  decoder := gob.NewDecoder(readfile)
  e = decoder.Decode(&decoded)

  if e != nil {
    fmt.Println("Error, the given file is not oat compatible")
    os.Exit(1)
  }

  readfile.Close()

  //convert the decoded value into json
  var jsondata, err = json.Marshal(decoded)

  if err != nil {
    fmt.Println("Error, given file cannot be read as oat")
    os.Exit(1)
  }

  //run the oat
  lang.OatRun(string(jsondata), params, dir)
}
