package main

import "os/exec"
import "strings"
import "encoding/json"

// #cgo CFLAGS: -std=c99
// #include "bind.h"
import "C"

func index(fileName, dir string, calcParams paramCalcOpts) {

  file := read("./pre.omm", "", true) + read(dir + fileName, "File Not Found: " + dir + fileName, true)
  fileNQ, _ := NQReplace(file)

  lexCmd := exec.Command("./lexer/main-win.exe")

  lexCmd.Stdin = strings.NewReader(fileNQ + "\n")

  _lex, _ := lexCmd.CombinedOutput()
  lex_ := string(_lex)

  var lex []string

  json.Unmarshal([]byte(lex_), &lex)

  var actions = actionizer(lex)

  var acts, _ = json.Marshal(actions)

  C.bind(C.CString(string(acts)), C.CString("HELLO"), C.CString("TES"))
}
