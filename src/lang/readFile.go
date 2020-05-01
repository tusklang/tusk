package main

import "fmt"
import "os"
import "strings"
import "encoding/json"
import "os/exec"

// #cgo CFLAGS: -std=c99
import "C"

func readFileJS(fileName string) []string {
  readCmd := exec.Command("./files/imports/index-win.exe")

  readCmd.Stdin = strings.NewReader(fileName)

  _file, _ := readCmd.CombinedOutput()
  file_ := strings.TrimSpace(string(_file))

  if strings.HasPrefix(file_, "Error") {
    fmt.Println(file_)
    os.Exit(1)
  }

  var files []string

  json.Unmarshal([]byte(file_), &files)

  return files
}
