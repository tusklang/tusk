package main

import "strings"
import "strconv"
import "reflect"

// #cgo CFLAGS: -std=c99
import "C"

type Condition struct {
  Type            string
  Condition     []Action
  Actions       []Action
}

type Action struct {
  Type            string
  Name            string
  ExpStr        []string
  ExpAct        []Action
  Params        []string
  Args        [][]Action
  Condition     []Condition
  ID              int

  //stuff for operations

  First         []Action
  Second        []Action
  Degree        []Action

  //stuff for indexes

  Value       [][]Action
  Indexes     [][]Action
  Hash_Values     map[string][]Action

  IsMutable       bool
}

var operations = []string{ "+", "-", "*", "/", "^", "%", "&", "|", "=", "!=", ">", "<", ">=", "<=", ")", "(", "~~", "~~~", ":" }

func convToAct(_val []interface{}, dir, name string) []Action {
  var val []Action

  if reflect.TypeOf(_val[0]).String() == "main.Lex" {

    var num []Lex

    for _, v := range _val {
      num = append(num, v.(Lex))
    }

    val = actionizer(num, true, dir, name)

  } else {

    for _, v := range _val {
      val = append(val, v.(Action))
    }

  }

  return val
}

func getLeft(index int, exp []interface{}, dir, name string) ([]Action, []interface{}) {

  var _num1 []interface{}

  //_num1 loop
  for o := index - 1; o >= 0; o-- {

    _num1 = append(_num1, exp[o])
  }

  reverseInterface(_num1)

  num1 := convToAct(_num1, dir, name)

  return num1, _num1
}

func getRight(index int, exp []interface{}, dir, name string) ([]Action, []interface{}) {
  var _num2 []interface{}

  //_num2 loop
  for o := index + 1; o < len(exp); o++ {

    _num2 = append(_num2, exp[o])
  }

  num2 := convToAct(_num2, dir, name)

  return num2, _num2
}

func calcExp(index int, exp []interface{}, dir, name string) ([]Action, []Action, []interface{}, []interface{}) {

  num1, _num1 := getLeft(index, exp, dir, name)
  num2, _num2 := getRight(index, exp, dir, name)

  return num1, num2, _num1, _num2
}

func callCalc(i *int, lex []Lex, len_lex int, dir, filename string) ([][]Action, string, [][]Action) {
  cbCnt := 0
  glCnt := 0
  bCnt := 0
  pCnt := 1

  var name = lex[(*i) + 2].Name

  indexes := [][]Lex{[]Lex{}}
  var putIndexes [][]Action

  if lex[(*i) + 3].Name == "." {

    cbCnt = 0
    glCnt = 0
    bCnt = 0
    pCnt = 0

    for o := (*i) + 4; o < len_lex; o++ {
      if lex[o].Name == "{" {
        cbCnt++
      }
      if lex[o].Name == "[:" {
        glCnt++
      }
      if lex[o].Name == "[" {
        bCnt++
      }
      if lex[o].Name == "(" {
        pCnt++
      }

      if lex[o].Name == "}" {
        cbCnt--
      }
      if lex[o].Name == ":]" {
        glCnt--
      }
      if lex[o].Name == "]" {
        bCnt--
      }
      if lex[o].Name == ")" {
        pCnt--
      }

      if lex[o].Name == "." {
        indexes = append(indexes, []Lex{})
      } else {

        (*i)++

        indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

        if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

          if o < len_lex - 1 && lex[o + 1].Name == "." {
            continue
          } else {
            break
          }

        }
      }
    }

    for _, v := range indexes {

      v = v[1:len(v) - 1]
      putIndexes = append(putIndexes, actionizer(v, true, dir, filename))
    }

    (*i)+=3
  }

  params := [][]Lex{[]Lex{}}

  for o := *i + 3; o < len_lex; o++ {
    if lex[o].Name == "{" {
      cbCnt++;
    }
    if lex[o].Name == "}" {
      cbCnt--;
    }

    if lex[o].Name == "[:" {
      glCnt++;
    }
    if lex[o].Name == ":]" {
      glCnt--;
    }

    if lex[o].Name == "[" {
      bCnt++;
    }
    if lex[o].Name == "]" {
      bCnt--;
    }

    if lex[o].Name == "(" {
      pCnt++;
    }
    if lex[o].Name == ")" {
      pCnt--;
    }

    if lex[o].Name == "(" {
      continue
    }

    if cbCnt != 0 && glCnt != 0 && bCnt != 0 && pCnt != 0 {
      params = append(params, []Lex{})
      continue
    }

    if lex[o].Name == ")" {
      break
    }

    if lex[o].Name == "," {
      params = append(params, []Lex{})
      continue
    }

    params[len(params) - 1] = append(params[len(params) - 1], lex[o])
  }

  var params_ = [][]Action{}

  for o := 0; o < len(params); o++ {

    if len(params[o]) == 0 {
      continue
    }

    params_ = append(params_, actionizer(params[o], true, dir, filename))
  }

  pCnt_ := 0
  skip_nums := 0

  for o := *i; o < len_lex; o++ {
    if lex[o].Name == "(" {
      pCnt_++
    }
    if lex[o].Name == ")" {
      pCnt_--
    }

    skip_nums++;

    if pCnt_ == 0 && lex[o].Name == "newlineS" {
      break
    }
  }

  (*i)+=skip_nums

  return params_, name, putIndexes
}

func procCalc(i *int, lex []Lex, len_lex int, dir, name string) ([]Action, []string, string) {

  var params []string
  var procName string
  var logic []Action

  if lex[(*i) + 1].Name == "~" {
    procName = lex[*i + 2].Name

    for o := (*i) + 4; o < len_lex; o++ {
      if lex[o].Name == ")" {
        break
      }

      if lex[o].Name == "," {
        (*i)++
        continue
      }

      params = append(params, lex[o].Name)
    }
    *i+=(len(params) + 5)

    var logic_ = []Lex{}

    cbCnt := 0

    for o := *i; o < len_lex; o++ {
      if lex[o].Name == "{" {
        cbCnt++
      }

      if lex[o].Name == "}" {
        cbCnt--
      }

      logic_ = append(logic_, lex[o])

      if cbCnt == 0 {
        break
      }
    }

    (*i)+=len(logic_) - 1

    logic = actionizer(logic_, false, dir, name)
  } else {
    params = []string{}
    procName = ""

    for o := (*i) + 2; o < len_lex; o+=2 {
      if lex[o].Name == ")" {
        break
      }

      params = append(params, lex[o].Name)
    }
    *i+=(3 + len(params))

    var logic_ = []Lex{}

    cbCnt := 0

    for o := *i; o < len_lex; o++ {
      if lex[o].Name == "{" {
        cbCnt++
      }

      if lex[o].Name == "}" {
        cbCnt--
      }

      logic_ = append(logic_, lex[o])

      if cbCnt == 0 {
        break
      }
    }

    (*i)+=len(logic_) - 1

    logic = actionizer(logic_, false, dir, name)
  }

  return logic, params, procName
}

func actionizer(lex []Lex, doExpress bool, dir, name string) []Action {
  var actions = []Action{}
  var len_lex = len(lex)

  for i := 0; i < len_lex; i++ {

    if doExpress {
      var exp []interface{}

      cbCnt := 0
      glCnt := 0
      bCnt := 0
      pCnt := 0

      for o := i; o < len_lex; o++ {
        if lex[o].Name == "{" {
          cbCnt++;
        }
        if lex[o].Name == "}" {
          cbCnt--;
        }

        if lex[o].Name == "[:" {
          glCnt++;
        }
        if lex[o].Name == ":]" {
          glCnt--;
        }

        if lex[o].Name == "[" {
          bCnt++;
        }
        if lex[o].Name == "]" {
          bCnt--;
        }

        if lex[o].Name == "(" {
          pCnt++;
        }
        if lex[o].Name == ")" {
          pCnt--;
        }

        if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
          break
        }

        exp = append(exp, lex[o])

        i++
      }

      for ;interfaceContainOperations(exp, "=") || interfaceContainOperations(exp, "!=") || interfaceContainOperations(exp, ">") || interfaceContainOperations(exp, "<") || interfaceContainOperations(exp, ">=") || interfaceContainOperations(exp, "<=") || interfaceContainOperations(exp, "~~") || interfaceContainOperations(exp, "~~~"); {
        indexes := map[string]int{
          "=": interfaceIndexOfOperations("=", exp),
          "!=": interfaceIndexOfOperations("!=", exp),
          ">": interfaceIndexOfOperations(">", exp),
          "<": interfaceIndexOfOperations("<", exp),
          ">=": interfaceIndexOfOperations(">=", exp),
          "<=": interfaceIndexOfOperations("<=", exp),
          "~~": interfaceIndexOfOperations("~~", exp),
          "~~~": interfaceIndexOfOperations("~~~", exp),
        }

        //get max index
        var min = [2]interface{}{}

        for k, v := range indexes {
          if v != -1 {
            min = [2]interface{}{ k, v }
          }
        }

        for k, v := range indexes {
          if (v != -1 && v > min[1].(int)) || min[1].(int) == -1 {
            min = [2]interface{}{ k, v }
          }
        }

        switch min[0].(string) {
          case "=":
            index := min[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "equals", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 47, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "!=":
            index := min[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "notEqual", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 48, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case ">":
            index := min[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "greater", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 49, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "<":
            index := min[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "less", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 50, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case ">=":
            index := min[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "greaterOrEqual", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 51, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "<=":
            index := min[1].(int)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "lessOrEqual", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 52, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "~~":
            index := min[1].(int)

            var degree_ []interface{}
            doDeg := false

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := index + 1; o < len(exp); o++ {
              if exp[o] == "{" {
                cbCnt++;
              }
              if exp[o] == "}" {
                cbCnt--;
              }

              if exp[o] == "[:" {
                glCnt++;
              }
              if exp[o] == ":]" {
                glCnt--;
              }

              if exp[o] == "[" {
                bCnt++;
              }
              if exp[o] == "]" {
                bCnt--;
              }

              if exp[o] == "(" {
                pCnt++;
              }
              if exp[o] == ")" {
                pCnt--;
              }

              if reflect.TypeOf(exp[o]).String() == "main.Lex" && exp[o].(Lex).Name == ":" && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
                doDeg = true
                break
              }

              degree_ = append(degree_, exp[o])
            }

            var degree = []Action{}
            var addDeg = 0

            if doDeg {
              degree = convToAct(degree_, dir, name)
              addDeg = len(degree_) + 1
            }

            num1, _num1 := getLeft(index, exp, dir, name)
            num2, _num2 := getRight(index + addDeg, exp, dir, name)

            var act_exp = Action{ "similar", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 54, num1, num2, degree, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + addDeg + 1:]...)

            exp = exp_
          case "~~~":
            index := min[1].(int)

            var degree_ []interface{}
            doDeg := false

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := index + 1; o < len(exp); o++ {
              if exp[o] == "{" {
                cbCnt++;
              }
              if exp[o] == "}" {
                cbCnt--;
              }

              if exp[o] == "[:" {
                glCnt++;
              }
              if exp[o] == ":]" {
                glCnt--;
              }

              if exp[o] == "[" {
                bCnt++;
              }
              if exp[o] == "]" {
                bCnt--;
              }

              if exp[o] == "(" {
                pCnt++;
              }
              if exp[o] == ")" {
                pCnt--;
              }

              if reflect.TypeOf(exp[o]).String() == "main.Lex" && exp[o].(Lex).Name == ":" && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
                doDeg = true
                break
              }

              degree_ = append(degree_, exp[o])
            }

            var degree = []Action{}
            var addDeg = 0

            if doDeg {
              degree = convToAct(degree_, dir, name)
              addDeg = len(degree_) + 1
            }

            num1, _num1 := getLeft(index, exp, dir, name)
            num2, _num2 := getRight(index + addDeg, exp, dir, name)

            var act_exp = Action{ "strictSimilar", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 55, num1, num2, degree, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + addDeg + 1:]...)

            exp = exp_
        }
      }

      for ;interfaceContainOperations(exp, "+") || interfaceContainOperations(exp, "-"); {

        if interfaceIndexOfOperations("+", exp) > interfaceIndexOfOperations("-", exp) || interfaceIndexOfOperations("-", exp) == -1 {
          index := interfaceIndexOfOperations("+", exp)

          num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

          var act_exp = Action{ "add", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 32, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

          exp_ := append(exp[:index - len(_num1)], act_exp)
          exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

          exp = exp_
        } else {
          index := interfaceIndexOfOperations("-", exp)

          num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

          var act_exp = Action{ "subtract", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 33, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

          exp_ := append(exp[:index - len(_num1)], act_exp)
          exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

          exp = exp_
        }

      }

      for ;interfaceContainOperations(exp, "*") || interfaceContainOperations(exp, "/") || interfaceContainOperations(exp, "%"); {

        indexes := map[string]int{
          "*": interfaceIndexOfOperations("*", exp),
          "/": interfaceIndexOfOperations("/", exp),
          "%": interfaceIndexOfOperations("%", exp),
        }

        //get max index
        var min = [2]interface{}{}

        for k, v := range indexes {
          if v != -1 {
            min = [2]interface{}{ k, v }
          }
        }

        for k, v := range indexes {
          if (v != -1 && v > min[1].(int)) || min[1].(int) == -1 {
            min = [2]interface{}{ k, v }
          }
        }

        switch min[0].(string) {
          case "*":
            index := interfaceIndexOfOperations("*", exp)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "multiply", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 34, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "/":
            index := interfaceIndexOfOperations("/", exp)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "divide", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 35, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
          case "%":
            index := interfaceIndexOfOperations("%", exp)

            num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

            var act_exp = Action{ "modulo", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 37, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

            exp_ := append(exp[:index - len(_num1)], act_exp)
            exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

            exp = exp_
        }

      }

      for ;interfaceContainOperations(exp, "^"); {
        index := interfaceIndexOfOperations("^", exp)

        num1, num2, _num1, _num2 := calcExp(index, exp, dir, name)

        var act_exp = Action{ "exponentiate", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 36, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

        exp_ := append(exp[:index - len(_num1)], act_exp)
        exp_ = append(exp_, exp[index + len(_num2) + 1:]...)

        exp = exp_
      }

      for ;interfaceContainOperations(exp, "!"); {

        index := interfaceIndexOfOperations("!", exp)

        var num []interface{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := index + 1; o < len(exp); o++ {

          if exp[o].(Lex).Name == "{" {
            cbCnt++
          }
          if exp[o].(Lex).Name == "}" {
            cbCnt--
          }

          if exp[o].(Lex).Name == "[" {
            bCnt++
          }
          if exp[o].(Lex).Name == "]" {
            bCnt--
          }

          if exp[o].(Lex).Name == "[:" {
            glCnt++
          }
          if exp[o].(Lex).Name == ":]" {
            glCnt--
          }

          if exp[o].(Lex).Name == "(" {
            pCnt++
          }
          if exp[o].(Lex).Name == ")" {
            pCnt--
          }

          if arrayContainInterface(operations, exp[o]) {
            break
          }

          num = append(num, exp[o])
        }

        numAct := convToAct(num, dir, name)

        var act_exp = Action{ "not", "operation", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 53, []Action{}, numAct, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

        exp_ := append(exp[:index], act_exp)
        exp_ = append(exp_, exp[index + len(num) + 1:]...)

        exp = exp_
      }

      var proc_indexes []int

      for ;interfaceContainWithProcIndex(exp, "(", proc_indexes); {

        index := interfaceIndexOfWithProcIndex("(", exp, proc_indexes)

        if index - 1 != -1 && (reflect.TypeOf(exp[index - 1]).String() != "main.Lex" || ((strings.HasPrefix(exp[index - 1].(Lex).Name, "$") || exp[index - 1].(Lex).Name == "]")))  {
          proc_indexes = append(proc_indexes, index)
          continue
        }

        var pExp []Lex

        pCnt := 0

        for o := index; o < len(exp); o++ {
          if exp[o].(Lex).Name == "(" {
            pCnt++;
          }
          if exp[o].(Lex).Name == ")" {
            pCnt--;
          }

          pExp = append(pExp, exp[o].(Lex))

          if pCnt == 0 {
            break
          }
        }

        pExp = pExp[1:len(pExp) - 1]

        pExpAct := actionizer(pExp, true, dir, name)

        scbCnt := 0
        sglCnt := 0
        sbCnt := 0
        spCnt := 0

        indexes := [][]Lex{}

        if !(index + len(pExp) + 2 >= len(exp)) {
          if exp[index + len(pExp) + 2].(Lex).Name == "." {
            for o := index + len(pExp) + 2; o < len_lex; o++ {
              if exp[o].(Lex).Name == "{" {
                scbCnt++
              }
              if exp[o].(Lex).Name == "}" {
                scbCnt--
              }

              if exp[o].(Lex).Name == "[" {
                sbCnt++
              }
              if exp[o].(Lex).Name == "]" {
                sbCnt--
              }

              if exp[o].(Lex).Name == "[:" {
                sglCnt++
              }
              if exp[o].(Lex).Name == ":]" {
                sglCnt--
              }

              if exp[o].(Lex).Name == "(" {
                spCnt++
              }
              if exp[o].(Lex).Name == ")" {
                spCnt--
              }

              if exp[o].(Lex).Name == "." {
                indexes = append(indexes, []Lex{})
              } else {

                i++

                indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], exp[o].(Lex))

                if scbCnt == 0 && sglCnt == 0 && sbCnt == 0 && spCnt == 0 {

                  if o < len(exp) - 1 && exp[o + 1].(Lex).Name == "." {
                    continue
                  } else {
                    break
                  }

                }
              }
            }

            var putIndexes [][]Action

            for _, v := range indexes {

              v = v[1:len(v) - 1]
              putIndexes = append(putIndexes, actionizer(v, true, dir, name))
            }

            pExpAct[0].Type = "expressionIndex"
            pExpAct[0].Indexes = putIndexes
            pExpAct[0].ID = 8 //set the action ID to the epxressionIndex ID
          }
        }

        exp = append([]interface{}{ pExpAct[0] }, exp...)
      }

      if reflect.TypeOf(exp[0]).String() == "main.Lex" {

        //variale that grets convved to a []Lex
        var toa []Lex

        for _, v := range exp {
          toa = append(toa, v.(Lex))
        }

        exp[0] = actionizer(toa, false, dir, name)[0]
      }

      actions = append(actions, exp[0].(Action))
    }

    if i >= len_lex {
      break
    }

    switch lex[i].Name {
      case "newlineN":
        actions = append(actions, Action{ "newline", "", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 0, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
      case "local":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 4; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o].Name == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, true, dir, name)

        actions = append(actions, Action{ "local", lex[i + 2].Name, []string{}, exp, []string{}, [][]Action{}, []Condition{}, 1, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(4 + len(exp_))
      case "dynamic":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 4; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o].Name == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, true, dir, name)

        actions = append(actions, Action{ "dynamic", lex[i + 2].Name, []string{}, exp, []string{}, [][]Action{}, []Condition{}, 2, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(4 + len(exp_))
      case "alt":

        var alter = Action{ "alt", "", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 3, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false }

        pCnt := 0

        cond_ := []Lex{}

        for o := i + 1; o < len_lex; o++ {
          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          cond_ = append(cond_, lex[o])

          if pCnt == 0 {
            break
          }
        }

        i+=len(cond_) + 1

        cond := actionizer(cond_, true, dir, name)

        for ;lex[i].Name == "=>"; {
          cbCnt := 0

          actions_ := []Lex{}

          for o := i + 1; o < len_lex; o++ {
            if lex[o].Name == "{" {
              cbCnt++
            }
            if lex[o].Name == "}" {
              cbCnt--
            }

            actions_ = append(actions_, lex[o])

            if cbCnt == 0 {
              break
            }
          }

          i+=len(actions_)
          actions := actionizer(actions_, true, dir, name)

          alter.Condition = append(alter.Condition, Condition{ "alt", cond, actions })
          i++

          if i >= len_lex {
            break
          }
        }

        actions = append(actions, alter)

      case "global":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 4; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o].Name == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, true, dir, name)

        actions = append(actions, Action{ "global", lex[i + 2].Name, []string{}, exp, []string{}, [][]Action{}, []Condition{}, 4, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(4 + len(exp_))
      case "log":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o].Name == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, true, dir, name)

        actions = append(actions, Action{ "log", "", []string{}, exp, []string{}, [][]Action{}, []Condition{}, 5, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(2 + len(exp_))
      case "print":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o].Name == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, false, dir, name)

        actions = append(actions, Action{ "print", "", []string{}, exp, []string{}, [][]Action{}, []Condition{}, 6, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(2 + len(exp_))
      case "{":
        exp_ := []Lex{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i; o < len_lex; o++ {

          if lex[o].Name == "{" {
            cbCnt++;
          }
          if lex[o].Name == "}" {
            cbCnt--;
          }

          if lex[o].Name == "[:" {
            glCnt++;
          }
          if lex[o].Name == ":]" {
            glCnt--;
          }

          if lex[o].Name == "[" {
            bCnt++;
          }
          if lex[o].Name == "]" {
            bCnt--;
          }

          if lex[o].Name == "(" {
            pCnt++;
          }
          if lex[o].Name == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          exp_ = append(exp_, lex[o])

          if cbCnt == 0 {
            break
          }

          if lex[o].Name == "newlineS" {
            break
          }
        }

        exp_ = exp_[1:]
        exp_ = exp_[:len(exp_) - 1]

        exp := actionizer(exp_, false, dir, name)

        actions = append(actions, Action{ "group", "", []string{}, exp, []string{}, [][]Action{}, []Condition{}, 9, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(len(exp_) + 1)
      case "process":

        putFalsey := make(map[string][]Action)
        putFalsey["falsey"] = []Action{ Action{ "falsey", "", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false } }

        logic, params, procName := procCalc(&i, lex, len_lex, dir, name)

        actions = append(actions, Action{ "process", procName, []string{}, logic, params, [][]Action{}, []Condition{}, 10, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, putFalsey, false })
      case "wait":

        var exp []Lex
        pCnt := 0

        for o := i + 1; o < len_lex; o++ {
          if lex[o].Name == "(" {
            pCnt++
            continue
          }
          if lex[o].Name == ")" {
            pCnt--
            continue
          }

          if pCnt == 0 {
            break
          }

          exp = append(exp, lex[o])
        }

        actionized := actionizer(exp, true, dir, name)

        actions = append(actions, Action{ "wait", "", []string{}, actionized, []string{}, [][]Action{}, []Condition{}, 57, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i++
      case "#":

        params_, name, putIndexes := callCalc(&i, lex, len_lex, dir, name)

        actions = append(actions, Action{ "#", name, []string{}, []Action{}, []string{}, params_, []Condition{}, 11, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), false })
      case "@":

        params_, name, putIndexes := callCalc(&i, lex, len_lex, dir, name)

        actions = append(actions, Action{ "@", name, []string{}, []Action{}, []string{}, params_, []Condition{}, 56, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), false })
      case "return":

        returner_ := []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
            break
          }

          returner_ = append(returner_, lex[o])
        }

        returner := actionizer(returner_, true, dir, name)

        actions = append(actions, Action{ "return", "", []string{}, returner, []string{}, [][]Action{}, []Condition{}, 12, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=len(returner_) + 2
      case "if":

        conditions := []Condition{}

        for ;lex[i].Name == "if"; {
          var cond_ = []Lex{}
          pCnt := 0

          for o := i + 1; o < len_lex; o++ {
            if lex[o].Name == "(" {
              pCnt++
            }
            if lex[o].Name == ")" {
              pCnt--
            }

            cond_ = append(cond_, lex[o])

            if pCnt == 0 {
              break
            }
          }

          cond := actionizer(cond_, true, dir, name)

          cbCnt := 0
          glCnt := 0
          bCnt := 0
          pCnt = 0

          actions_ := []Lex{}

          var curlyBraceCond = lex[i + 1 + len(cond_)].Name == "{"

          for o := i + 1 + len(cond_); o < len_lex; o++ {
            if lex[o].Name == "{" {
              cbCnt++
            }
            if lex[o].Name == "}" {
              cbCnt--
            }

            if lex[o].Name == "[:" {
              glCnt++
            }
            if lex[o].Name == ":]" {
              glCnt--
            }

            if lex[o].Name == "[" {
              bCnt++
            }
            if lex[o].Name == "]" {
              bCnt--
            }

            if lex[o].Name == "(" {
              pCnt++
            }
            if lex[o].Name == ")" {
              pCnt--
            }

            actions_ = append(actions_, lex[o])

            if !curlyBraceCond {

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                break
              }
            } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
              break
            }
          }

          acts := actionizer(actions_, false, dir, name)

          var if_ = Condition{ "if", cond, acts }

          conditions = append(conditions, if_)

          i+=(1 + len(cond_) + len(actions_))

          if i >= len_lex {
            break
          }
        }

        if !(i >= len_lex) {
          for ;lex[i].Name == "elseif"; {

            var cond_ = []Lex{}
            pCnt := 0

            for o := i + 1; o < len_lex; o++ {
              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              cond_ = append(cond_, lex[o])

              if pCnt == 0 {
                break
              }
            }

            cond := actionizer(cond_, true, dir, name)

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt = 0

            actions_ := []Lex{}

            var curlyBraceCond = lex[i + 1 + len(cond_)].Name == "{"

            for o := i + 1 + len(cond_); o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              actions_ = append(actions_, lex[o])

              if !curlyBraceCond {

                if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                  break
                }
              } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
                break
              }
            }

            acts := actionizer(actions_, false, dir, name)

            var elseif_ = Condition{ "elseif", cond, acts }

            conditions = append(conditions, elseif_)

            i+=(1 + len(cond_) + len(actions_))

            if i >= len_lex {
              break
            }
          }
        }

        if !(i >= len_lex) {
          actions_ := []Lex{}

          for ;lex[i].Name == "else"; {
            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            //allow for user to write: else ~ <do something>;
            var curlyBraceCond = lex[i + 1].Name == "{"

            for o := i + 1; o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if !curlyBraceCond {

                if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                  break
                }
              } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
                break
              }

              actions_ = append(actions_, lex[o])
            }

            if !curlyBraceCond {
              actions_ = actions_[1:]
            }

            actions := actionizer(actions_, false, dir, name)

            var else_ = Condition{ "else", []Action{}, actions }

            conditions = append(conditions, else_)

            i+=(1 + len(actions_))

            if i >= len_lex {
              break
            }
          }
        }

        actions = append(actions, Action{ "conditional", "", []string{}, []Action{}, []string{}, [][]Action{}, conditions, 13, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i--
      case "import":

        var fileDir = lex[i + 2].Name

        //remove the quotes
        fileDir = fileDir[1:len(fileDir) - 1]

        var files []string

        //see if user wants to import a fire from the basedir
        if strings.HasPrefix(fileDir, "?~") {

          if strings.HasPrefix(fileDir[2:], "/") {
            files = readFileJS("./stdlib" + fileDir[2:])
          } else {
            files = readFileJS("./stdlib/" + fileDir[2:])
          }

        } else {
          files = readFileJS(dir + fileDir)
        }

        var lexxed [][]Lex

        for _, v := range files {
          lexxed = append(lexxed, lexer(v, dir, name))
        }

        var actionizedFiles [][]Action

        for _, v := range lexxed {
          actionizedFiles = append(actionizedFiles, actionizer(v, false, dir, name))
        }

        actions = append(actions, Action{ "import", "", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 14, []Action{}, []Action{}, []Action{}, actionizedFiles, [][]Action{}, make(map[string][]Action), false })
        i+=3
      case "read":
        var phrase = []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionizedPhrase := actionizer(phrase, true, dir, name)
        actions = append(actions, Action{ "read", "", []string{}, actionizedPhrase, []string{}, [][]Action{}, []Condition{}, 15, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(2 + len(phrase))
      case "break":
        actions = append(actions, Action{ "break", "", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 16, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
      case "skip":
        actions = append(actions, Action{ "skip", "", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 17, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
      case "eval":
        var phrase = []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionized := actionizer(phrase, true, dir, name)
        actions = append(actions, Action{ "eval", "", []string{}, actionized, []string{}, [][]Action{}, []Condition{}, 18, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(2 + len(phrase))
      case "typeof":
        var phrase = []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionizedPhrase := actionizer(phrase, true, dir, name)
        actions = append(actions, Action{ "typeof", "", []string{}, actionizedPhrase, []string{}, [][]Action{}, []Condition{}, 19, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(2 + len(phrase))
      case "loop":

        var condition_ = []Lex{}

        pCnt := 0

        for o := i + 1; o < len_lex; o++ {
          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          condition_ = append(condition_, lex[o])

          if pCnt == 0 {
            break
          }
        }

        condition := actionizer(condition_, true, dir, name)
        action_ := []Lex{}

        var curlyBraceCond = lex[i + 1 + len(condition_)].Name == "{"

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt = 0

        for o := i + 1 + len(condition_); o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          action_ = append(action_, lex[o])

          if !curlyBraceCond {

            if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
              break
            }
          } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
            break
          }
        }

        action := actionizer(action_, false, dir, name)

        actions = append(actions, Action{ "loop", "", []string{}, action, []string{}, [][]Action{}, []Condition{ { "loop", condition, action } }, 21, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(1 + len(condition_) + len(action_))
      case "[:":
        var phrase = []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          phrase = append(phrase, lex[o])

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
            break
          }
        }

        i+=len(phrase)

        phrase = phrase[1:len(phrase) - 1]

        var _translated =  [][][]Lex{ [][]Lex{ []Lex{}, []Lex{} } }

        cbCnt = 0
        glCnt = 0
        bCnt = 0
        pCnt = 0

        cur := 0

        for _, v := range phrase {

          if v.Name == "newlineN" {
            continue
          }

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

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && v.Name == ":" {
            cur = 1
            continue
          }
          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && v.Name == "," {
            cur = 0
            _translated = append(_translated, [][]Lex{ []Lex{}, []Lex{} })
            continue
          }

          _translated[len(_translated) - 1][cur] = append(_translated[len(_translated) - 1][cur], v)
        }

        var translated = make(map[string][]Action)

        translated["falsey"] = []Action{ Action{ "falsey", "exp_value", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false } }

        for _, v := range _translated {

          if len(v[0]) <= 0 {
            break
          }

          var name_ = v[0][0].Name

          if strings.HasPrefix(v[0][0].Name, "'") {
            name_ = name_[1:len(name_) - 1]
          }

          translated[name_] = actionizer(v[1], true, dir, name)
        }

        if i >= len_lex {
          actions = append(actions, Action{ "hash", "hashed_value", []string{""}, []Action{}, []string{}, [][]Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, translated, false })
          break
        }

        isMutable := false

        //checks for a runtime hash
        for ;i < len_lex && lex[i].Name == ":::"; i++ {
          isMutable = !isMutable
        }

        if i >= len_lex {
          actions = append(actions, Action{ "hash", "hashed_value", []string{""}, []Action{}, []string{}, [][]Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, translated, isMutable })
          break
        }

        if lex[i].Name == "." {
          indexes := [][]Lex{ []Lex{} }

          cbCnt = 0
          glCnt = 0
          bCnt = 0
          pCnt = 0

          for o := i + 1; o < len_lex; o++ {
            if lex[o].Name == "{" {
              cbCnt++
            }
            if lex[o].Name == "[:" {
              glCnt++
            }
            if lex[o].Name == "[" {
              bCnt++
            }
            if lex[o].Name == "(" {
              pCnt++
            }

            if lex[o].Name == "}" {
              cbCnt--
            }
            if lex[o].Name == ":]" {
              glCnt--
            }
            if lex[o].Name == "]" {
              bCnt--
            }
            if lex[o].Name == ")" {
              pCnt--
            }

            if lex[o].Name == "." {
              indexes = append(indexes, []Lex{})
            } else {

              i++

              indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                if o < len_lex - 1 && lex[o + 1].Name == "." {
                  continue
                } else {
                  break
                }

              }
            }
          }

          var putIndexes [][]Action

          for _, v := range indexes {

            v = v[1:len(v) - 1]
            putIndexes = append(putIndexes, actionizer(v, true, dir, name))
          }

          i+=3

          actions = append(actions, Action{ "hashIndex", "", []string{""}, []Action{}, []string{}, [][]Action{}, []Condition{}, 23, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, translated, isMutable })
        } else {
          actions = append(actions, Action{ "hash", "hashed_value", []string{""}, []Action{}, []string{}, [][]Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, translated, isMutable })
        }
      case "[":
        var phrase = []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          phrase = append(phrase, lex[o])

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
            break
          }
        }

        i+=len(phrase)

        phrase = phrase[1:len(phrase) - 1]

        var arr [][]Action

        for o := 0; o < len(phrase); o++ {

          var sub []Lex

          cbCnt := 0
          glCnt := 0
          bCnt := 0
          pCnt := 0

          for j := o; j < len(phrase); j++ {

            if phrase[j].Name == "{" {
              cbCnt++
            }
            if phrase[j].Name == "}" {
              cbCnt--
            }

            if phrase[j].Name == "[:" {
              glCnt++
            }
            if phrase[j].Name == ":]" {
              glCnt--
            }

            if phrase[j].Name == "[" {
              bCnt++
            }
            if phrase[j].Name == "]" {
              bCnt--
            }

            if phrase[j].Name == "(" {
              pCnt++
            }
            if phrase[j].Name == ")" {
              pCnt--
            }

            if phrase[j].Name == "," && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
              break
            }
            sub = append(sub, phrase[j])
          }

          o+=len(sub)

          arr = append(arr, actionizer(sub, true, dir, name))
        }

        hashedArr := make(map[string][]Action)

        hashedArr["falsey"] = []Action{ Action{ "falsey", "exp_value", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false } }

        cur := "0"

        for _, v := range arr {
          hashedArr[cur] = v
          cur = add(cur, "1", paramCalcOpts{}, -1)
        }

        if i >= len_lex {
          actions = append(actions, Action{ "array", "hashed_value", []string{""}, []Action{}, []string{}, [][]Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, hashedArr, false })
          break
        }

        isMutable := false

        //checks for a runtime array
        for ;i < len_lex && lex[i].Name == ":::"; i++ {
          isMutable = !isMutable
        }

        if i >= len_lex {
          actions = append(actions, Action{ "array", "hashed_value", []string{""}, []Action{}, []string{}, [][]Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, hashedArr, isMutable })
          break
        }

        if lex[i].Name == "." {
          indexes := [][]Lex{ []Lex{} }

          cbCnt = 0
          glCnt = 0
          bCnt = 0
          pCnt = 0

          for o := i + 1; o < len_lex; o++ {
            if lex[o].Name == "{" {
              cbCnt++
            }
            if lex[o].Name == "[:" {
              glCnt++
            }
            if lex[o].Name == "[" {
              bCnt++
            }
            if lex[o].Name == "(" {
              pCnt++
            }

            if lex[o].Name == "}" {
              cbCnt--
            }
            if lex[o].Name == ":]" {
              glCnt--
            }
            if lex[o].Name == "]" {
              bCnt--
            }
            if lex[o].Name == ")" {
              pCnt--
            }

            if lex[o].Name == "." {
              indexes = append(indexes, []Lex{})
            } else {

              i++

              indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                if o < len_lex - 1 && lex[o + 1].Name == "." {
                  continue
                } else {
                  break
                }

              }
            }
          }

          var putIndexes [][]Action

          for _, v := range indexes {

            v = v[1:len(v) - 1]
            putIndexes = append(putIndexes, actionizer(v, true, dir, name))
          }

          i+=3

          actions = append(actions, Action{ "arrayIndex", "", []string{""}, []Action{}, []string{}, [][]Action{}, []Condition{}, 25, []Action{}, []Action{}, []Action{}, arr, putIndexes, hashedArr, isMutable })
        } else {
          actions = append(actions, Action{ "array", "hashed_value", []string{""}, []Action{}, []string{}, [][]Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, hashedArr, isMutable })
        }
      case "ascii":
        var phrase = []Lex{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionizedPhrase := actionizer(phrase, true, dir, name)
        actions = append(actions, Action{ "ascii", "", []string{}, actionizedPhrase, []string{}, [][]Action{}, []Condition{}, 26, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=(2 + len(phrase))
      case "len":

        var exp []Lex
        pCnt := 0

        for o := i + 1; o < len_lex; o++ {
          if lex[o].Name == "(" {
            pCnt++
            continue
          }
          if lex[o].Name == ")" {
            pCnt--
            continue
          }

          if pCnt == 0 {
            break
          }

          exp = append(exp, lex[o])
        }

        actionized := actionizer(exp, true, dir, name)

        actions = append(actions, Action{ "len", "", []string{}, actionized, []string{}, [][]Action{}, []Condition{}, 31, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        i+=3 + len(exp)
      case "each":
        var condition_ = []Lex{}

        pCnt := 0

        for o := i + 1; o < len_lex; o++ {
          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          condition_ = append(condition_, lex[o])

          if pCnt == 0 {
            break
          }
        }

        i+=len(condition_) + 1

        condition_ = condition_[1:len(condition_) - 1]

        var _iterator []Lex

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt = 0

        var stopIterIndex int

        for k, v := range condition_ {

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

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && v.Name == "newlineS" {

            //index where the iterator stopped
            stopIterIndex = k
            break
          }

          _iterator = append(_iterator, v)
          stopIterIndex = k
        }

        iterator := actionizer(_iterator, true, dir, name)

        var var1 string
        var var2 string

        if stopIterIndex + 1 >= len(condition_) {
          var1 = "$k"
          var2 = "$v"
        } else if stopIterIndex + 3 >= len(condition_) {
          var1 = condition_[stopIterIndex + 1].Name
          var2 = "$v"
        } else {
          var1, var2 = condition_[stopIterIndex + 1].Name, condition_[stopIterIndex + 3].Name
        }

        cbCnt = 0
        glCnt = 0
        bCnt = 0
        pCnt = 0

        var exp []Lex

        var curlyBraceCond = lex[i].Name == "{"

        for o := i; o < len_lex; o++ {
          if lex[o].Name == "{" {
            cbCnt++
          }
          if lex[o].Name == "}" {
            cbCnt--
          }

          if lex[o].Name == "[:" {
            glCnt++
          }
          if lex[o].Name == ":]" {
            glCnt--
          }

          if lex[o].Name == "[" {
            bCnt++
          }
          if lex[o].Name == "]" {
            bCnt--
          }

          if lex[o].Name == "(" {
            pCnt++
          }
          if lex[o].Name == ")" {
            pCnt--
          }

          exp = append(exp, lex[o])

          if !curlyBraceCond {

            if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
              break
            }
          } else if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
            break
          }
        }

        i+=len(exp) + 1
        actionized := actionizer(exp, false, dir, name)
        actions = append(actions, Action{ "each", "", []string{ var1, var2 }, actionized, []string{}, [][]Action{}, []Condition{}, 59, iterator, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })

      //file system keywords

      case "files.read":

        args := cproc(&i, lex, uint(1), "files.read", dir, name)
        actions = append(actions, Action{ "files.read", "", []string{}, []Action{}, []string{}, args, []Condition{}, 60, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })

      case "files.write":

        args := cproc(&i, lex, uint(2), "files.write", dir, name)
        actions = append(actions, Action{ "files.write", "", []string{}, []Action{}, []string{}, args, []Condition{}, 61, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })

      case "files.exists":

        args := cproc(&i, lex, uint(1), "files.exists", dir, name)
        actions = append(actions, Action{ "files.exists", "", []string{}, []Action{}, []string{}, args, []Condition{}, 62, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })

      case "files.isFile":

        args := cproc(&i, lex, uint(1), "files.isFile", dir, name)
        actions = append(actions, Action{ "files.isFile", "", []string{}, []Action{}, []string{}, args, []Condition{}, 63, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })

      case "files.isDir":

        args := cproc(&i, lex, uint(1), "files.isDir", dir, name)
        actions = append(actions, Action{ "files.isDir", "", []string{}, []Action{}, []string{}, args, []Condition{}, 64, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
      //////////////////////

      case "kill":

        if lex[i + 1].Name == "<-" {
          actions = append(actions, Action{ "kill_thread", "", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 65, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
          i++
        } else {
          actions = append(actions, Action{ "kill", "", []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, 66, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
        }

      default:

        valPuts := func(lex []Lex, i int) int {

          if i >= len_lex {
            return 1
          }

          isMutable := false
          val := lex[i].Name

          //checks for a runtime value
          for ;i + 1 < len_lex && lex[i + 1].Name == ":::"; i++ {
            isMutable = !isMutable
          }

          i++

          switch C.GoString(GetType(C.CString(val))) {

            case "string": {

              noQ := val[1:len(val) - 1]
              hashedString := make(map[string][]Action)

              //specify the value for the "falsey" case
              hashedString["falsey"] = []Action{ Action{ "falsey", "exp_value", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), isMutable } }

              cur := "0"

              for _, v := range noQ {

                hashedIndex := make(map[string][]Action)
                hashedIndex["falsey"] = []Action{ Action{ "falsey", "exp_value", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), isMutable } }

                hashedString[cur] = []Action{ Action{ "string", "exp_value", []string{ string(v) }, []Action{}, []string{}, [][]Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashedIndex, isMutable } }
                cur = add(cur, "1", paramCalcOpts{}, -1)
              }

              actions = append(actions, Action{ "string", "exp_value", []string{ noQ }, []Action{}, []string{}, [][]Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashedString, isMutable })
            }
            case "number":

              hashed := make(map[string][]Action)

              //specify the value for the "falsey" case
              hashed["falsey"] = []Action{ Action{ "falsey", "exp_value", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), isMutable } }

              actions = append(actions, Action{ "number", "exp_value", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, false })
            case "boolean":

              hashed := make(map[string][]Action)

              //specify the value for the "falsey" case
              hashed["falsey"] = []Action{ Action{ "boolean", "exp_value", []string{ strconv.FormatBool(val != "true") }, []Action{}, []string{}, [][]Action{}, []Condition{}, 40, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), isMutable } }

              actions = append(actions, Action{ "boolean", "exp_value", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 40, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, isMutable })
            case "falsey":

              hashed := make(map[string][]Action)

              //specify the value for the "falsey" case
              hashed["falsey"] = []Action{ Action{ "falsey", "exp_value", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), isMutable } }

              actions = append(actions, Action{ "falsey", "exp_value", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, isMutable })
            case "none":

              if strings.HasPrefix(val, "$") {

                actions = append(actions, Action{ "variable", val, []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 43, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), isMutable })
              } else {

                hashedString := make(map[string][]Action)

                //specify the value for the "falsey" case
                hashedString["falsey"] = []Action{ Action{ "falsey", "exp_value", []string{ "undef" }, []Action{}, []string{}, [][]Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), isMutable } }

                //get it? 42?
                actions = append(actions, Action{ "none", "exp_value", []string{ val }, []Action{}, []string{}, [][]Action{}, []Condition{}, 42, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), isMutable })
              }
          }

          return 0
        }

        if i + 1 < len_lex {

          if lex[i + 1].Name == "->" {

            var val_ []Lex

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 2; o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && arrayContainInterface(operations, lex[o].Name) {
                break
              }

              val_ = append(val_, lex[o])
            }

            val := actionizer(val_, true, dir, name)

            getValueType := func(val string) string {
              switch (C.GoString(GetType(C.CString(val)))) {
                case "string": fallthrough
                case "number": fallthrough
                case "boolean": fallthrough
                case "falsey":
                  return "exp_value"
                case "array": fallthrough
                case "hash":
                  return "hashed_value"
              }

              return "exp_value"
            }

            actions = append(actions, Action{ "cast", lex[i].Name, []string{ getValueType(lex[i].Name) }, val, []string{}, [][]Action{}, []Condition{}, 58, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
            i+=len(val_) + 2
            continue
          }

          if (lex[i + 1].Name == "++" || lex[i + 1].Name == "--") && strings.HasPrefix(lex[i].Name, "$") {

            id_ := []byte(lex[i + 1].Name)

            id := ""

            for o := 0; o < len(id_); o++ {
              _id := strconv.Itoa(int(id_[o]))
              id+=_id
            }

            intID, _ := strconv.Atoi(id)

            actions = append(actions, Action{ lex[i + 1].Name, lex[i].Name, []string{}, []Action{}, []string{}, [][]Action{}, []Condition{}, intID, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
            i++
            continue;
          }

          if (lex[i + 1].Name == "+=" || lex[i + 1].Name == "-=" || lex[i + 1].Name == "*=" || lex[i + 1].Name == "/=" || lex[i + 1].Name == "%=" || lex[i + 1].Name == "^=") && strings.HasPrefix(lex[i].Name, "$") {

            var by_ []Lex

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 2; o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                break
              }

              by_ = append(by_, lex[o])
            }

            by := actionizer(by_, true, dir, name)

            id_ := []byte(lex[i + 1].Name)

            id := ""

            for o := 0; o < len(id_); o++ {
              _id := strconv.Itoa(int(id_[o]))
              id+=_id
            }

            intID, _ := strconv.Atoi(id)

            actions = append(actions, Action{ lex[i + 1].Name, lex[i].Name, []string{}, by, []string{}, [][]Action{}, []Condition{}, intID, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), false })
            continue;
          }

          doPutIndex := false

          icbCnt := 0
          iglCnt := 0
          ibCnt := 0
          ipCnt := 0

          for o := i; o < len_lex; o++ {
            if lex[o].Name == "{" {
              icbCnt++
            }
            if lex[o].Name == "}" {
              icbCnt--
            }

            if lex[o].Name == "[:" {
              iglCnt++
            }
            if lex[o].Name == ":]" {
              iglCnt--
            }

            if lex[o].Name == "[" {
              ibCnt++
            }
            if lex[o].Name == "]" {
              ibCnt--
            }

            if lex[o].Name == "(" {
              ipCnt++
            }
            if lex[o].Name == ")" {
              ipCnt--
            }

            if icbCnt == 0 && iglCnt == 0 && ibCnt == 0 && ipCnt == 0 && lex[o].Name == "newlineS" {
              break
            }

            if icbCnt == 0 && iglCnt == 0 && ibCnt == 0 && ipCnt == 0 && lex[o].Name == ":" {
              doPutIndex = true
              break
            }
          }

          var indexes [][]Action
          varname := lex[i].Name

          if lex[i + 1].Name == "." && lex[i + 2].Name == "[" && doPutIndex {

            _indexes := [][]Lex{}

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 1; o < len_lex; i, o = i + 1, o + 1 {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if lex[o].Name == "." {
                _indexes = append(_indexes, []Lex{})
                continue
              }

              _indexes[len(_indexes) - 1] = append(_indexes[len(_indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o + 1].Name == ":" {
                break
              }
            }

            for _, v := range _indexes {
              indexes = append(indexes, actionizer(v[1:len(v) - 1], true, dir, name))
            }

            i++
          }

          if lex[i + 1].Name == ":" && (strings.HasPrefix(lex[i].Name, "$") || lex[i].Name == "]") {
            exp_ := []Lex{}

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 2; o < len_lex; o++ {

              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "}" {
                cbCnt--
              }

              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == ":]" {
                glCnt--
              }

              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "]" {
                bCnt--
              }

              if lex[o].Name == "(" {
                pCnt++
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o].Name == "newlineS" {
                break
              }

              exp_ = append(exp_, lex[o]);
            }

            exp := actionizer(exp_, true, dir, name)

            actions = append(actions, Action{ "let", varname, []string{}, exp, []string{}, [][]Action{}, []Condition{}, 28, []Action{}, []Action{}, []Action{}, [][]Action{}, indexes, make(map[string][]Action), false })
            i+=(len(exp))
            continue
          }

          if lex[i + 1].Name == "." {

            val := lex[i].Name

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            indexes := [][]Lex{ []Lex{} }

            cbCnt = 0
            glCnt = 0
            bCnt = 0
            pCnt = 0

            for o := i + 2; o < len_lex; o++ {
              if lex[o].Name == "{" {
                cbCnt++
              }
              if lex[o].Name == "[:" {
                glCnt++
              }
              if lex[o].Name == "[" {
                bCnt++
              }
              if lex[o].Name == "(" {
                pCnt++
              }

              if lex[o].Name == "}" {
                cbCnt--
              }
              if lex[o].Name == ":]" {
                glCnt--
              }
              if lex[o].Name == "]" {
                bCnt--
              }
              if lex[o].Name == ")" {
                pCnt--
              }

              if lex[o].Name == "." {
                indexes = append(indexes, []Lex{})
              } else {

                i++

                indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

                if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                  if o < len_lex - 1 && lex[o + 1].Name == "." {
                    continue
                  } else {
                    break
                  }

                }
              }
            }

            var putIndexes [][]Action

            for _, v := range indexes {

              v = v[1:len(v) - 1]
              putIndexes = append(putIndexes, actionizer(v, true, dir, name))
            }

            i+=3

            if strings.HasPrefix(val, "$") {
              actVal := actionizer([]Lex{ Lex{ val, "", 0, "", "", dir } }, true, dir, name)

              actions = append(actions, Action{ "variableIndex", "", []string{}, actVal, []string{}, [][]Action{}, []Condition{}, 46, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), false })
            } else {
              actVal := actionizer([]Lex{ Lex{ val, "", 0, "", "", dir } }, true, dir, name)

              actions = append(actions, Action{ "expressionIndex", "", []string{}, actVal, []string{}, [][]Action{}, []Condition{}, 8, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), false })
            }

          }

          valPuts(lex, i)

        } else {

          valPuts(lex, i)
        }
      }
  }

  return actions
}
