package lang

import "strings"
import "reflect"

func arrayContain(arr []string, sub string) bool {

  for i := 0; i < len(arr); i++ {
    if arr[i] == sub {
      return true
    }
  }

  return false;
}

//the list of cprocs goes here
//just add to the slice if you add a new cproc
var CPROCS = []string{ "files.read", "files.write", "files.exists", "files.isFile", "files.isDir", "regex.match", "regex.replace", "this" }

func arrayContainInterface(arr []string, sub interface{}) bool {

  for i := 0; i < len(arr); i++ {
    if arr[i] == sub {
      return true
    }
  }

  return false;
}

func arrayContainInterfaceLexStr(arr []string, sub interface{}) bool {

  loop:
  for i := 0; i < len(arr); i++ {

    switch sub.(type) {
      case Action: continue loop
    }

    if arr[i] == sub.(Lex).Name {
      return true
    }
  }

  return false;
}

func arrayContainInterfaceOperations(arr []string, sub interface{}) bool {

  loop:
  for i := 0; i < len(arr); i++ {

    switch sub.(type) {
      case Lex: continue loop
      case Action: continue loop
    }

    if arr[i] == sub.(string) {
      return true
    }
  }

  return false;
}

func arrayContain2Nest(arr [][]string, sub string) bool {

  for i := 0; i < len(arr); i++ {
    if arrayContain(arr[i], sub) {
      return true
    }
  }

  return false
}

func indexOf2Nest(sub string, arr [][]string) []int {
  for i := 0; i < len(arr); i++ {
    for o := 0; o < len(arr[i]); o++ {
      if arr[i][o] == sub {
        return []int{ i, o }
      }
    }
  }

  return []int{ -1, -1 }
}

func RepeatAdd(s string, times int) string {
  returner := ""

  for ;times > 0; times-- {
    returner+=s
  }

  return returner
}

func indexOf(sub string, data []string) int {
  for k, v := range data {
    if sub == v {
      return k
    }
  }
  return -1
}

func interfaceContain(inter []interface{}, sub interface{}) bool {
  for _, a := range inter {
    if a == sub {
      return true
    }
  }
  return false
}

func interfaceIndexOf(sub interface{}, inter []interface{}) int {
  for k, v := range inter {
    if sub == v {
      return k
    }
  }
  return -1
}

func interfaceContainOperations(inter []interface{}, sub interface{}) bool {

  cbCnt := 0
  glCnt := 0
  bCnt := 0
  pCnt := 0

  loop:
  for _, a := range inter {

    switch a.(type) {
      case Action:
        continue loop
    }

    if a.(Lex).Name == "{" {
      cbCnt++;
    }
    if a.(Lex).Name == "}" {
      cbCnt--;
    }

    if a.(Lex).Name == "[:" {
      glCnt++;
    }
    if a.(Lex).Name == ":]" {
      glCnt--;
    }

    if a.(Lex).Name == "[" {
      bCnt++;
    }
    if a.(Lex).Name == "]" {
      bCnt--;
    }

    if a.(Lex).Name == "(" {
      pCnt++;
    }
    if a.(Lex).Name == ")" {
      pCnt--;
    }

    if a.(Lex).Name == sub && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
      return true
    }
  }

  return false
}

func interfaceIndexOfOperations(sub interface{}, inter []interface{}) int {

  cbCnt := 0
  glCnt := 0
  bCnt := 0
  pCnt := 0

  loop:
  for k, a := range inter {

    switch a.(type) {
      case Action:
        continue loop
    }

    if a.(Lex).Name == "{" {
      cbCnt++;
    }
    if a.(Lex).Name == "}" {
      cbCnt--;
    }

    if a.(Lex).Name == "[:" {
      glCnt++;
    }
    if a.(Lex).Name == ":]" {
      glCnt--;
    }

    if a.(Lex).Name == "[" {
      bCnt++;
    }
    if a.(Lex).Name == "]" {
      bCnt--;
    }

    if a.(Lex).Name == "(" {
      pCnt++;
    }
    if a.(Lex).Name == ")" {
      pCnt--;
    }

    if a.(Lex).Name == sub && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
      return k
    }
  }

  return -1
}

func interfaceContainWithProcIndex(inter []interface{}, sub interface{}, indexes []int) bool {

  loop:
  for k, v := range inter {

    switch v.(type) {
      case Action:
        continue loop
    }

    if k != 0 {

      //only if inter[k - 1] is a lex
      switch inter[k - 1].(type) {
        case Lex:

          //if inter[k - 1] is a process or a variable, continue the loop
          if arrayContain(CPROCS, inter[k - 1].(Lex).Name) || inter[k - 1].(Lex).Name == "process" || strings.HasPrefix(inter[k - 1].(Lex).Name, "$") || inter[k - 1].(Lex).Name == "]" {
            continue loop
          }
      }
    }

    if sub.(string) == v.(Lex).Name {

      for _, o := range indexes {
        if k == o {
          continue loop
        }

      }

      return true
    }
  }
  return false
}

func interfaceIndexOfWithProcIndex(sub interface{}, inter []interface{}, indexes []int) int {

  loop:
  for k, v := range inter {

    switch v.(type) {
      case Action:
        continue loop
    }

    if k != 0 {

      //only if inter[k - 1] is a lex
      switch inter[k - 1].(type) {
        case Lex:

          //if inter[k - 1] is a cproc continue the loop
          if arrayContain(CPROCS, inter[k - 1].(Lex).Name) {
            continue loop
          }
      }
    }

    if sub.(string) == v.(Lex).Name {

      for _, o := range indexes {
        if k == o {
          continue loop
        }

      }

      return k
    }
  }
  return -1
}

func interfaceContainForExp(inter []interface{}, _sub []string) bool {

  var sub []interface{}

  for _, v := range _sub {
    sub = append(sub, v)
  }

  cbCnt := 0
  glCnt := 0
  bCnt := 0
  pCnt := 0

  for o := 0; o < len(inter); o++ {

    v := inter[o]

    if reflect.TypeOf(v).String() != "main.Lex" {
      continue
    }

    //prevent parenthesis after process declarations from being counted as expression parenthesis
    if o > 0 && !(strings.HasPrefix(inter[o - 1].(Lex).Name, "$") || inter[o - 1].(Lex).Name == "]" || inter[o - 1].(Lex).Name == "process" || arrayContain(CPROCS, inter[o - 1].(Lex).Name) ) && v.(Lex).Name == "(" {

      scbCnt := 0
      sglCnt := 0
      sbCnt := 0
      spCnt := 0

      for i := o; i < len(inter); i, o = i + 1, o + 1 {
        if inter[i].(Lex).Name == "{" {
          scbCnt++;
        }
        if inter[i].(Lex).Name == "}" {
          scbCnt--;
        }

        if inter[i].(Lex).Name == "[:" {
          sglCnt++;
        }
        if inter[i].(Lex).Name == ":]" {
          sglCnt--;
        }

        if inter[i].(Lex).Name == "[" {
          sbCnt++;
        }
        if inter[i].(Lex).Name == "]" {
          sbCnt--;
        }

        if inter[i].(Lex).Name == "(" {
          spCnt++;
        }
        if inter[i].(Lex).Name == ")" {
          spCnt--;
        }

        if scbCnt == 0 && sglCnt == 0 && sbCnt == 0 && spCnt == 0 {
          break
        }
      }

      continue
    }

    if v.(Lex).Name == "{" {
      cbCnt++;
    }
    if v.(Lex).Name == "}" {
      cbCnt--;
    }

    if v.(Lex).Name == "[:" {
      glCnt++;
    }
    if v.(Lex).Name == ":]" {
      glCnt--;
    }

    if v.(Lex).Name == "[" {
      bCnt++;
    }
    if v.(Lex).Name == "]" {
      bCnt--;
    }

    if v.(Lex).Name == "(" {
      pCnt++;
    }
    if v.(Lex).Name == ")" {
      pCnt--;
    }

    if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

      for _, i := range sub {

        if i == "(" {
          if o > 0 && (strings.HasPrefix(inter[o - 1].(Lex).Name, "$") || inter[o - 1].(Lex).Name == "]" || inter[o - 1].(Lex).Name == "process") {
            continue
          }
        }

        if i == v.(Lex).Name {
          return true
        }
      }

    }
  }

  return false
}
