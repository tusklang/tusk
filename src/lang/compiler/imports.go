package compiler

import "fmt"
import "os"
import "strings"
import "encoding/json"
import "os/exec"

//export ReadFileJS
func ReadFileJS(fileName string) []map[string]string {
  readCmd := exec.Command("./lang/compiler/imports/index-win.exe")

  readCmd.Stdin = strings.NewReader(fileName)

  _file, _ := readCmd.CombinedOutput()
  file_ := string(_file)

  if strings.HasPrefix(file_, "Error") {
    fmt.Println(file_)
    os.Exit(1)
  }

  var files []map[string]string

  json.Unmarshal([]byte(file_), &files)

  return files
}
