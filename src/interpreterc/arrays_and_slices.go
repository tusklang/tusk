package main

import "strings"

func arrayContain(arr []string, sub string) bool {

  for i := 0; i < len(arr); i++ {
    if arr[i] == sub {
      return true
    }
  }

  return false;
}

func arrayContainInterface(arr []string, sub interface{}) bool {

  for i := 0; i < len(arr); i++ {
    if arr[i] == sub {
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

func interfaceContainWithProcIndex(inter []interface{}, sub interface{}, indexes []int) bool {

  loop:
  for k, v := range inter {
    if sub == v {

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
    if sub == v {

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

    if o > 0 && (strings.HasPrefix(inter[o - 1].(string), "$") || inter[o - 1].(string) == "]") && v == "(" {

      scbCnt := 0
      sglCnt := 0
      sbCnt := 0
      spCnt := 0

      for i := o; i < len(inter); i, o = i + 1, o + 1 {
        if inter[i] == "{" {
          scbCnt++;
        }
        if inter[i] == "}" {
          scbCnt--;
        }

        if inter[i] == "[:" {
          sglCnt++;
        }
        if inter[i] == ":]" {
          sglCnt--;
        }

        if inter[i] == "[" {
          sbCnt++;
        }
        if inter[i] == "]" {
          sbCnt--;
        }

        if inter[i] == "(" {
          spCnt++;
        }
        if inter[i] == ")" {
          spCnt--;
        }

        if scbCnt == 0 && sglCnt == 0 && sbCnt == 0 && spCnt == 0 {
          break
        }
      }

      continue
    }

    if v == "{" {
      cbCnt++;
    }
    if v == "}" {
      cbCnt--;
    }

    if v == "[:" {
      glCnt++;
    }
    if v == ":]" {
      glCnt--;
    }

    if v == "[" {
      bCnt++;
    }
    if v == "]" {
      bCnt--;
    }

    if v == "(" {
      pCnt++;
    }
    if v == ")" {
      pCnt--;
    }

    if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

      for _, i := range sub {
        if i == v {
          return true
        }
      }

    }
  }

  return false
}
