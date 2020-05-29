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
var imported = []string{}

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

func getType(val string) string {

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
    return "string"
  } else if strings.HasPrefix(val, "[:") {
    return "hash"
  } else if strings.HasPrefix(val, "[") {
    return "array"
  } else if val == "true" || val == "false" {
    return "boolean"
  } else if val == "undef" || val == "null" {
    return "falsey"
  } else if numMatch(val) {
    return "number"
  }

  return "none"
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
    C.colorprint(C.CString("Error while lexxing in " + fmt.Sprint(v.(map[string]interface{})["Dir"]) + "!"), C.int(12))
    fmt.Print(fmt.Sprint(v.(map[string]interface{})["Error"]), "\n\n" + strings.Repeat("-", 90) + "\n")
    fmt.Println()
  }

  for _, v := range ret["WARNS"] {
    C.colorprint(C.CString("Warning while lexxing in " + fmt.Sprint(v.(map[string]interface{})["Dir"]) + "! "), C.int(14))
    fmt.Print(fmt.Sprint(v.(map[string]interface{})["Error"]), "\n\n" + strings.Repeat("-", 90) + "\n")
    fmt.Println()
  }

  if len(ret["ERRORS"]) != 0 {
    fmt.Println("Script was lexxed with", len(ret["ERRORS"]), "errors, and", len(ret["WARNS"]), "warnings")
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
func OatRun(acts, cli_params, dir string) {

  argv := make([]*C.char, len(os.Args[1:]))

  for k, v := range os.Args[1:] {
    cstring := C.CString(v)
    argv[k] = cstring
  }

  _ = argv

  C.bindParser(C.CString(acts), C.CString(cli_params), C.CString(dir), C.int(len(os.Args[1:])), &argv[0])
}

//export Run
func Run(params map[string]map[string]interface{}) {

  dir := params["Files"]["DIR"].(string)
  fileName := params["Files"]["NAME"].(string)

  imported = append(imported, dir + fileName)

  file := ReadFileJS(dir + fileName)[0]["Content"]

  lex := Lexer(file, dir, fileName)

  var actions = Actionizer(lex, false, dir, fileName)

  cp, _ := json.Marshal(params)
  acts, _ := json.MarshalIndent(actions, "", "  ")

  argv := make([]*C.char, len(os.Args[1:]))

  for k, v := range os.Args[1:] {
    cstring := C.CString(v)
    argv[k] = cstring
  }

  C.bindParser(C.CString(string(acts)), C.CString(string(cp)), C.CString(dir), C.int(len(os.Args[1:])), &argv[0])
}
