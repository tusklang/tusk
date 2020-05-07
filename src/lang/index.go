package lang

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

var operators = []string{"^", "*", "/", "%", "+", "-", "&", "|", "!", "~", ";"}

//export Kill
func Kill() {
  os.Exit(1)
}

type Lex struct {
  Name   string
  Exp    string
  Line   uint64
  Type   string
  OName  string
  Dir    string
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
  } else if val == "undef" || val == "null" {
    return C.CString("falsey")
  } else if numMatch(val) {
    return C.CString("number")
  }

  return C.CString("none")
}

//export Lexer
func Lexer(file, dir, name string) []Lex {
  lexCmd := exec.Command("./lang/lexer/main-win.exe")

  var in = map[string]string{
    "f": file,
    "dir": dir,
    "name": name,
  }

  //dont know why, but json.Marshal does not work, so you I used json.MarshalIndent
  instr, _ := json.MarshalIndent(in, "", "  ")

  lexCmd.Stdin = strings.NewReader(string(instr))

  _out, _ := lexCmd.CombinedOutput()
  out := string(_out)

  var ret map[string][]interface{}

  json.Unmarshal([]byte(out), &ret)

  for _, v := range ret["ERRORS"] {
    C.colorprint(C.CString("Error while lexxing in " + dir + name + "!"), C.int(12))
    fmt.Print(fmt.Sprint(v), "\n\n" + strings.Repeat("-", 90) + "\n")
    fmt.Println()
  }

  for _, v := range ret["WARNS"] {
    C.colorprint(C.CString("Warning while lexxing in " + dir + name + "! "), C.int(14))
    fmt.Print(fmt.Sprint(v), "\n\n" + strings.Repeat("-", 90) + "\n")
    fmt.Println()
  }

  if len(ret["ERRORS"]) != 0 {
    os.Exit(1)
  }

  var lex []Lex

  for _, v := range ret["LEX"] {

    encoded, _ := json.Marshal(v)

    var cur Lex

    json.Unmarshal([]byte(encoded), &cur)

    lex = append(lex, cur)
  }

  return lex
}

//export OatRun
func OatRun(acts, cli_params string) {
  C.bindCgo(C.CString(acts), C.CString(cli_params))
}

//export Run
func Run(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"].(string)
  fileName := params["Files"]["NAME"].(string)

  file := ReadFileJS(dir + fileName)[0]["Content"]

  lex := Lexer(file, dir, fileName)

  var actions = Actionizer(lex, false, dir, fileName)

  var acts, _ = json.Marshal(actions)

  cp, _ := json.Marshal(params)

  _, _ = acts, cp

  C.bindCgo(C.CString(string(acts)), C.CString(string(cp)))
}
