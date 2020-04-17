package main

import "os"
import "os/exec"
import "strings"
import "encoding/json"
import "unicode"

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

  acts, _ := json.Marshal(actionizer(lex, false))

  return C.CString(string(acts))
}

//export GetType
func GetType(cVal *C.char) *C.char {

  val := C.GoString(cVal)

  var numMatch = func(num string) bool {
    for _, v := range num {
      if !unicode.IsDigit(v) && v != '.' {
        return false
      }
    }

    return true
  }

  if strings.HasPrefix(val, "\"") || strings.HasPrefix(val, "'") || strings.HasPrefix(val, "`") {
    return C.CString("string")
  } else if strings.HasPrefix(val, "[:") {
    return C.CString("hash")
  } else if strings.HasPrefix(val, "[") {
    return C.CString("array")
  } else if val == "true" || val == "false" {
    return C.CString("boolean")
  } else if val == "undefined" || val == "null" {
    return C.CString("falsey")
  } else if numMatch(val) {
    return C.CString("number")
  }

  return C.CString("none")
}

//export CLex
func CLex(_file *C.char) *C.char {

  file := C.GoString(_file)

  lexCmd := exec.Command("./lexer/main-win.exe")

  fileNQ, _ := NQReplace(file)

  lexCmd.Stdin = strings.NewReader(fileNQ + "\n")

  _lex, _ := lexCmd.CombinedOutput()
  lex_ := string(_lex)

  return C.CString(lex_)
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

  var actions = actionizer(lex, false)

  var acts, _ = json.Marshal(actions)

  cp, _ := json.Marshal(calcParams)

  C.bind(C.CString(string(acts)), C.CString(string(cp)), C.CString(dir))
}
