package main

import "os/exec"
import "strings"
import "encoding/json"

// #cgo CXXFLAGS: -I/usr/lib/
// #cgo LDFLAGS: -L/usr/lib/ -lstdc++
//#include <parser/bridge.hpp>
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

  C.bridge()
}
