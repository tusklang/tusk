package run

import "encoding/json"
import "strings"
import "fmt"
import "os"

import "lang" //omm language

//export Run
func Run(params map[string]map[string]interface{}) {
  dir := params["Files"]["DIR"].(string)
  fileName := params["Files"]["NAME"].(string)

  //read the oat
  file := strings.TrimSpace(lang.ReadFileJS(dir + fileName)[0])

  paramsJ, _ := json.Marshal(params)

  //determine if the given file is an oat
  var testoat map[string]interface{}
  err := json.Unmarshal([]byte(file), &testoat)

  _ = testoat

  //if there is an error display that there is an error with the oat
  if err == nil {
    fmt.Println("Error, the given file is not in an oat format")
    os.Exit(1)
  }
  ///////////////////////////////////////

  lang.OatRun(file, string(paramsJ))
}
