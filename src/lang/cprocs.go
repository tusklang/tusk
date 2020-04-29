package main

import "fmt"
import "os"
import "strings"

//this is just a function to actionize c processes in omm
func cproc(i *int, lex []Lex, PARAM_COUNT uint, name, dir string) [][]Action {

  var curLex = lex[(*i)]
  var paramExp []Lex

  pCnt := 0

  //get what is in the parenthesis
  for (*i)++; (*i) < len(lex); (*i)++ {

    if lex[(*i)].Name == "(" {
      pCnt++
    }
    if lex[(*i)].Name == ")" {
      pCnt--
    }

    paramExp = append(paramExp, lex[(*i)])

    if pCnt == 0 {
      break;
    }
  }

  //remove the parenthesis
  paramExp = paramExp[1:len(paramExp) - 1]

  cbCnt := 0
  glCnt := 0
  bCnt := 0
  pCnt = 0

  var splitParams = [][]Lex{ []Lex{} }

  for _, v := range paramExp {

    if v.Name == "{" {
      cbCnt++
    }
    if v.Name == "}" {
      cbCnt--
    }

    if v.Name == "[:" {
      glCnt++
    }
    if v.Name == ":]" {
      glCnt--
    }

    if v.Name == "[" {
      bCnt++
    }
    if v.Name == "]" {
      bCnt--
    }

    if v.Name == "(" {
      pCnt++
    }
    if v.Name == ")" {
      pCnt--
    }

    if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && v.Name == "," {
      splitParams = append(splitParams, []Lex{})
      continue
    }

    splitParams[len(splitParams) - 1] = append(splitParams[len(splitParams) - 1], v)

  }

  var actionSplit [][]Action

  for _, v := range splitParams {
    actionSplit = append(actionSplit, actionizer(v, true, dir))
  }

  if uint(len(splitParams)) < PARAM_COUNT {

    //throw an error
    fmt.Println("Error while actionizing! Not enough arguments for", /* say the process */ name, "\n\nError occured on line", curLex.Line, "\nFound near:", strings.TrimSpace(curLex.Exp))

    //exit the process
    os.Exit(1)
  }

  return actionSplit
}
