package main

import "os"
import "os/exec"
import "strings"
import "encoding/json"
import "unicode"
import "regexp"
import "fmt"

// #cgo CFLAGS: -std=c99
// #include "bind.h"
import "C"

//export Kill
func Kill() {
  os.Exit(1)
}

type Lex struct {
  Name   string
  Exp  []string
  Line   uint64
}

//export Cactions
func Cactions(file *C.char, dir *C.char) *C.char {

  var lex []Lex

  json.Unmarshal([]byte(C.GoString(file)), &lex)

  acts, _ := json.Marshal(actionizer(lex, false, C.GoString(dir)))

  return C.CString(string(acts))
}

//export GetType
func GetType(cVal *C.char) *C.char {

  val := C.GoString(cVal)

  var numMatch = func(num string) bool {

    //see if it includes at least one digit
    match, _ := regexp.MatchString("\\d", num)

    if !match {
      return false
    }

    for _, v := range num {
      if !unicode.IsDigit(v) && v != '.' && v != '-' && v != '+' {
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

func lexer(file string) []Lex {
  fileNQ, _ := NQReplace(file)

  lexCmd := exec.Command("./lexer/main-win.exe")

  lexCmd.Stdin = strings.NewReader(fileNQ + "\n")

  _lex, _ := lexCmd.CombinedOutput()
  lex_ := string(_lex)

  if strings.HasPrefix(lex_, "Error") {
    fmt.Println(lex_)
    os.Exit(1)
  }

  var lex []Lex

  json.Unmarshal([]byte(lex_), &lex)

  return lex
}

func index(fileName, dir string, calcParams paramCalcOpts) {

  file := read(dir + fileName, "File Not Found: " + dir + fileName, true)

  lex := lexer(file)

  var actions = actionizer(lex, false, dir)

  var acts, _ = json.Marshal(actions)

  cp, _ := json.Marshal(calcParams)

  _, _ = acts, cp

  C.bindCgo(C.CString(string(acts)), C.CString(string(cp)), C.CString(dir))
}
