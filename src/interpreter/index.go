package main

import "os/exec"
import "strings"
import "encoding/json"

type Funcs struct {
  Name string
  Line uint64
}

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

  parser(actions, calcParams, dir, 1, []Funcs{ Funcs{ file, 0 } }, make(map[string]Variable), false)
}
