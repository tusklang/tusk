package main

import "os"
import "os/exec"
import "strings"
import "encoding/json"
import "regexp"

// #cgo CFLAGS: -std=c99
// #include "bind.h"
import "C"

//export Kill
func Kill() {
  os.Exit(1)
}

//export Cactions
func Cactions(file *C.char) *C.char {

  var lex []string

  json.Unmarshal([]byte(C.GoString(file)), &lex)

  acts, _ := json.Marshal(actionizer(lex))

  return C.CString(string(acts))
}

//export GetType
func GetType(cVal *C.char) *C.char {
  var val [][]string

  json.Unmarshal([]byte(C.GoString(cVal)), &val);

  numMatch, _ := regexp.MatchString("^(\\d|\\.)", val[0][0])

  if strings.HasPrefix(val[0][0], "\"") || strings.HasPrefix(val[0][0], "'") || strings.HasPrefix(val[0][0], "`") {
    return C.CString("string")
  } else if strings.HasPrefix(val[0][0], "[:") {
    return C.CString("hash")
  } else if strings.HasPrefix(val[0][0], "[") {
    return C.CString("array")
  } else if val[0][0] == "true" || val[0][0] == "false" {
    return C.CString("boolean")
  } else if val[0][0] == "undefined" || val[0][0] == "null" {
    return C.CString("falsey")
  } else if numMatch {
    return C.CString("number")
  }

  return C.CString("none")
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

   var acts, _ = json.Marshal(actions)

  cp, _ := json.Marshal(calcParams)

  C.bind(C.CString(string(acts)), C.CString(string(cp)), C.CString(dir))
}
