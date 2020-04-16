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
  Args          []Action
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

  //type of the action as a value
  ValueType     []Action
  IsMutable       bool
}

func actionizer(lex []string, doExpress bool) []Action {
  var actions = []Action{}
  var len_lex = len(lex)

  var operations = []string{ "+", "-", "*", "/", "^", "%", "&", "|", "=", ">", "<", ">=", "<=", ")", "(", "~~", "~~~", "!" }

  for i := 0; i < len_lex; i++ {

    if doExpress {
      var exp []interface{}

      cbCnt := 0
      glCnt := 0
      bCnt := 0
      pCnt := 0

      for o := i; o < len_lex; o++ {
        if lex[o] == "{" {
          cbCnt++;
        }
        if lex[o] == "}" {
          cbCnt--;
        }

        if lex[o] == "[:" {
          glCnt++;
        }
        if lex[o] == ":]" {
          glCnt--;
        }

        if lex[o] == "[" {
          bCnt++;
        }
        if lex[o] == "]" {
          bCnt--;
        }

        if lex[o] == "(" {
          pCnt++;
        }
        if lex[o] == ")" {
          pCnt--;
        }

        if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
          break
        }

        exp = append(exp, lex[o])

        i++
      }

      if !interfaceContainForExp(exp, operations) {

        var act_exp []string

        for _, v := range exp {
          act_exp = append(act_exp, v.(string))
        }

        return actionizer(act_exp, false);
      }

      var proc_indexes []int

      for ;interfaceContainWithProcIndex(exp, "(", proc_indexes); {

        index := interfaceIndexOfWithProcIndex("(", exp, proc_indexes)

        if index - 1 != -1 && (strings.HasPrefix(exp[index - 1].(string), "$") || exp[index - 1].(string) == "len")  {
          proc_indexes = append(proc_indexes, index)
          continue
        }

        var pExp []string

        pCnt := 0

        for o := index; o < len(exp); o++ {
          if exp[o] == "(" {
            pCnt++;
          }
          if exp[o] == ")" {
            pCnt--;
          }

          pExp = append(pExp, exp[o].(string))

          if pCnt == 0 {
            break
          }
        }

        pExp = pExp[1:len(pExp) - 1]

        pExpAct := actionizer(pExp, true)

        exp_ := append(exp[:index], pExpAct[0])
        exp_ = append(exp_, exp[index + len(pExp):]...)
        exp = exp_
      }

      for ;interfaceContain(exp, "^"); {
        index := interfaceIndexOf("^", exp)

        var _num1 = []interface{}{}
        var _num2 = []interface{}{}

        //_num1 loop
        for o := index - 1; o >= 0; o-- {

          if arrayContainInterface(operations, exp[o]) {
            break
          }

          _num1 = append(_num1, exp[o])
        }

        //_num2 loop
        for o := index + 1; o < len(exp); o++ {

          if arrayContainInterface(operations, exp[o]) {
            break
          }

          _num2 = append(_num2, exp[o])
        }

        var num1 []Action
        var num2 []Action

        if reflect.TypeOf(_num1[0]).String() == "string" {

          var num []string

          for _, v := range _num1 {
            num = append(num, v.(string))
          }

          num1 = actionizer(num, true)

        } else {

          for _, v := range _num1 {
            num1 = append(num1, v.(Action))
          }

        }

        if reflect.TypeOf(_num2[0]).String() == "string" {

          var num []string

          for _, v := range _num2 {
            num = append(num, v.(string))
          }

          num2 = actionizer(num, true)

        } else {

          for _, v := range _num2 {
            num2 = append(num2, v.(Action))
          }

        }

        var act_exp = Action{ "exponentiate", "operation", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 36, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�operation" }, false), false }

        exp_ := append(exp[:index - 1], act_exp)
        exp_ = append(exp_, exp[index + 1:])

        exp = exp_
      }

      for ;interfaceContain(exp, "*") || interfaceContain(exp, "/"); {

        if interfaceIndexOf("*", exp) > interfaceIndexOf("/", exp) || interfaceIndexOf("/", exp) == -1 {
          index := interfaceIndexOf("*", exp)

          var _num1 = []interface{}{}
          var _num2 = []interface{}{}

          //_num1 loop
          for o := index - 1; o >= 0; o-- {

            if arrayContainInterface(operations, exp[o]) {
              break
            }

            _num1 = append(_num1, exp[o])
          }

          //_num2 loop
          for o := index + 1; o < len(exp); o++ {

            if arrayContainInterface(operations, exp[o]) {
              break
            }

            _num2 = append(_num2, exp[o])
          }

          var num1 []Action
          var num2 []Action

          if reflect.TypeOf(_num1[0]).String() == "string" {

            var num []string

            for _, v := range _num1 {
              num = append(num, v.(string))
            }

            num1 = actionizer(num, false)

          } else {

            for _, v := range _num1 {
              num1 = append(num1, v.(Action))
            }

          }

          if reflect.TypeOf(_num2[0]).String() == "string" {

            var num []string

            for _, v := range _num2 {
              num = append(num, v.(string))
            }

            num2 = actionizer(num, false)

          } else {

            for _, v := range _num2 {
              num2 = append(num2, v.(Action))
            }

          }

          var act_exp = Action{ "multiply", "operation", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 34, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�operation" }, false), false }

          exp_ := append(exp[:index - 1], act_exp)
          exp_ = append(exp_, exp[index + 1:])

          exp = exp_
        } else {
          index := interfaceIndexOf("/", exp)

          var _num1 = []interface{}{}
          var _num2 = []interface{}{}

          //_num1 loop
          for o := index - 1; o >= 0; o-- {

            if arrayContainInterface(operations, exp[o]) {
              break
            }

            _num1 = append(_num1, exp[o])
          }

          //_num2 loop
          for o := index + 1; o < len(exp); o++ {

            if arrayContainInterface(operations, exp[o]) {
              break
            }

            _num2 = append(_num2, exp[o])
          }

          var num1 []Action
          var num2 []Action

          if reflect.TypeOf(_num1[0]).String() == "string" {

            var num []string

            for _, v := range _num1 {
              num = append(num, v.(string))
            }

            num1 = actionizer(num, false)

          } else {

            for _, v := range _num1 {
              num1 = append(num1, v.(Action))
            }

          }

          if reflect.TypeOf(_num2[0]).String() == "string" {

            var num []string

            for _, v := range _num2 {
              num = append(num, v.(string))
            }

            num2 = actionizer(num, false)

          } else {

            for _, v := range _num2 {
              num2 = append(num2, v.(Action))
            }

          }

          var act_exp = Action{ "divide", "operation", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 35, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�operation" }, false), false }

          exp_ := append(exp[:index - 1], act_exp)
          exp_ = append(exp_, exp[index + 1:])

          exp = exp_
        }

      }

      for ;interfaceContain(exp, "%"); {
        index := interfaceIndexOf("%", exp)

        var _num1 = []interface{}{}
        var _num2 = []interface{}{}

        //_num1 loop
        for o := index - 1; o >= 0; o-- {

          if arrayContainInterface(operations, exp[o]) {
            break
          }

          _num1 = append(_num1, exp[o])
        }

        //_num2 loop
        for o := index + 1; o < len(exp); o++ {

          if arrayContainInterface(operations, exp[o]) {
            break
          }

          _num2 = append(_num2, exp[o])
        }

        var num1 []Action
        var num2 []Action

        if reflect.TypeOf(_num1[0]).String() == "string" {

          var num []string

          for _, v := range _num1 {
            num = append(num, v.(string))
          }

          num1 = actionizer(num, false)

        } else {

          for _, v := range _num1 {
            num1 = append(num1, v.(Action))
          }

        }

        if reflect.TypeOf(_num2[0]).String() == "string" {

          var num []string

          for _, v := range _num2 {
            num = append(num, v.(string))
          }

          num2 = actionizer(num, false)

        } else {

          for _, v := range _num2 {
            num2 = append(num2, v.(Action))
          }

        }

        var act_exp = Action{ "modulo", "operation", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 37, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�operation" }, false), false }

        exp_ := append(exp[:index - 1], act_exp)
        exp_ = append(exp_, exp[index + 1:])

        exp = exp_
      }

      for ;interfaceContain(exp, "+") || interfaceContain(exp, "-"); {

        if interfaceIndexOf("+", exp) > interfaceIndexOf("-", exp) || interfaceIndexOf("-", exp) == -1 {
          index := interfaceIndexOf("+", exp)

          var _num1 []interface{}
          var _num2 []interface{}

          cbCnt := 0
          glCnt := 0
          bCnt := 0
          pCnt := 0

          //_num1 loop
          for o := index - 1; o >= 0; o-- {

            if lex[o] == "{" {
              cbCnt++;
            }
            if lex[o] == "}" {
              cbCnt--;
            }

            if lex[o] == "[:" {
              glCnt++;
            }
            if lex[o] == ":]" {
              glCnt--;
            }

            if lex[o] == "[" {
              bCnt++;
            }
            if lex[o] == "]" {
              bCnt--;
            }

            if lex[o] == "(" {
              pCnt++;
            }
            if lex[o] == ")" {
              pCnt--;
            }

            if arrayContainInterface(operations, exp[o]) && !(cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0) {
              break
            }

            _num1 = append(_num1, exp[o])
          }

          //_num2 loop
          for o := index + 1; o < len(exp); o++ {

            if lex[o] == "{" {
              cbCnt++;
            }
            if lex[o] == "}" {
              cbCnt--;
            }

            if lex[o] == "[:" {
              glCnt++;
            }
            if lex[o] == ":]" {
              glCnt--;
            }

            if lex[o] == "[" {
              bCnt++;
            }
            if lex[o] == "]" {
              bCnt--;
            }

            if lex[o] == "(" {
              pCnt++;
            }
            if lex[o] == ")" {
              pCnt--;
            }

            if arrayContainInterface(operations, exp[o]) && !(cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0) {
              break
            }

            _num2 = append(_num2, exp[o])
          }

          reverseInterface(_num1)

          var num1 []Action
          var num2 []Action

          if reflect.TypeOf(_num1[0]).String() == "string" {

            var num []string

            for _, v := range _num1 {
              num = append(num, v.(string))
            }

            num1 = actionizer(num, false)

          } else {

            for _, v := range _num1 {
              num1 = append(num1, v.(Action))
            }

          }

          if reflect.TypeOf(_num2[0]).String() == "string" {

            var num []string

            for _, v := range _num2 {
              num = append(num, v.(string))
            }

            num2 = actionizer(num, true)

          } else {

            for _, v := range _num2 {
              num2 = append(num2, v.(Action))
            }

          }

          var act_exp = Action{ "add", "operation", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 32, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�operation" }, false), false }

          exp_ := append(exp[:index - len(_num1)], act_exp)
          exp_ = append(exp_, exp[index + len(_num2) + 1:])

          exp = exp_
        } else {
          index := interfaceIndexOf("-", exp)

          var _num1 []interface{}
          var _num2 []interface{}

          //_num1 loop
          for o := index - 1; o >= 0; o-- {

            if arrayContainInterface(operations, exp[o]) {
              break
            }

            _num1 = append(_num1, exp[o])
          }

          //_num2 loop
          for o := index + 1; o < len(exp); o++ {

            if arrayContainInterface(operations, exp[o]) {
              break
            }

            _num2 = append(_num2, exp[o])
          }

          var num1 []Action
          var num2 []Action

          if reflect.TypeOf(_num1[0]).String() == "string" {

            var num []string

            for _, v := range _num1 {
              num = append(num, v.(string))
            }

            num1 = actionizer(num, false)

          } else {

            for _, v := range _num1 {
              num1 = append(num1, v.(Action))
            }

          }

          if reflect.TypeOf(_num2[0]).String() == "string" {

            var num []string

            for _, v := range _num2 {
              num = append(num, v.(string))
            }

            num2 = actionizer(num, false)

          } else {

            for _, v := range _num2 {
              num2 = append(num2, v.(Action))
            }

          }

          var act_exp = Action{ "subtract", "operation", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 33, num1, num2, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�operation" }, false), false }

          exp_ := append(exp[:index - 1], act_exp)
          exp_ = append(exp_, exp[index + 1:])

          exp = exp_
        }

      }

      if reflect.TypeOf(exp[0]).String() == "string" {
        exp[0] = actionizer([]string{ exp[0].(string) }, false)[0]
      }

      actions = append(actions, exp[0].(Action))
    }

    if i >= len_lex {
      break
    }

    switch lex[i] {
      case "newlineN":
        actions = append(actions, Action{ "newline", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 0, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
      case "local":
        exp_ := []string{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 4; o < len_lex; o++ {

          if lex[o] == "{" {
            cbCnt++;
          }
          if lex[o] == "}" {
            cbCnt--;
          }

          if lex[o] == "[:" {
            glCnt++;
          }
          if lex[o] == ":]" {
            glCnt--;
          }

          if lex[o] == "[" {
            bCnt++;
          }
          if lex[o] == "]" {
            bCnt--;
          }

          if lex[o] == "(" {
            pCnt++;
          }
          if lex[o] == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o] == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, true)

        actions = append(actions, Action{ "local", lex[i + 2], []string{}, exp, []string{}, []Action{}, []Condition{}, 1, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(4 + len(exp_))
      case "dynamic":
        exp_ := []string{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 4; o < len_lex; o++ {

          if lex[o] == "{" {
            cbCnt++;
          }
          if lex[o] == "}" {
            cbCnt--;
          }

          if lex[o] == "[:" {
            glCnt++;
          }
          if lex[o] == ":]" {
            glCnt--;
          }

          if lex[o] == "[" {
            bCnt++;
          }
          if lex[o] == "]" {
            bCnt--;
          }

          if lex[o] == "(" {
            pCnt++;
          }
          if lex[o] == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o] == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, true)

        actions = append(actions, Action{ "dynamic", lex[i + 2], []string{}, exp, []string{}, []Action{}, []Condition{}, 2, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(4 + len(exp_))
      case "alt":

        var alter = Action{ "alt", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 3, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false }

        pCnt := 0

        cond_ := []string{}

        for o := i + 1; o < len_lex; o++ {
          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          cond_ = append(cond_, lex[o])

          if pCnt == 0 {
            break
          }
        }

        i+=len(cond_) + 1

        cond := actionizer(cond_, true)

        for ;lex[i] == "=>"; {
          cbCnt := 0

          actions_ := []string{}

          for o := i + 1; o < len_lex; o++ {
            if lex[o] == "{" {
              cbCnt++
            }
            if lex[o] == "}" {
              cbCnt--
            }

            actions_ = append(actions_, lex[o])

            if cbCnt == 0 {
              break
            }
          }

          i+=len(actions_)
          actions := actionizer(actions_, true)

          alter.Condition = append(alter.Condition, Condition{ "alt", cond, actions })
          i++

          if i >= len_lex {
            break
          }
        }

        actions = append(actions, alter)

      case "global":
        exp_ := []string{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 4; o < len_lex; o++ {

          if lex[o] == "{" {
            cbCnt++;
          }
          if lex[o] == "}" {
            cbCnt--;
          }

          if lex[o] == "[:" {
            glCnt++;
          }
          if lex[o] == ":]" {
            glCnt--;
          }

          if lex[o] == "[" {
            bCnt++;
          }
          if lex[o] == "]" {
            bCnt--;
          }

          if lex[o] == "(" {
            pCnt++;
          }
          if lex[o] == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o] == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, true)

        actions = append(actions, Action{ "global", lex[i + 2], []string{}, exp, []string{}, []Action{}, []Condition{}, 4, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(4 + len(exp_))
      case "log":
        exp_ := []string{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {

          if lex[o] == "{" {
            cbCnt++;
          }
          if lex[o] == "}" {
            cbCnt--;
          }

          if lex[o] == "[:" {
            glCnt++;
          }
          if lex[o] == ":]" {
            glCnt--;
          }

          if lex[o] == "[" {
            bCnt++;
          }
          if lex[o] == "]" {
            bCnt--;
          }

          if lex[o] == "(" {
            pCnt++;
          }
          if lex[o] == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o] == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, true)

        actions = append(actions, Action{ "log", "", exp_, exp, []string{}, []Action{}, []Condition{}, 5, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(2 + len(exp_))
      case "print":
        exp_ := []string{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {

          if lex[o] == "{" {
            cbCnt++;
          }
          if lex[o] == "}" {
            cbCnt--;
          }

          if lex[o] == "[:" {
            glCnt++;
          }
          if lex[o] == ":]" {
            glCnt--;
          }

          if lex[o] == "[" {
            bCnt++;
          }
          if lex[o] == "]" {
            bCnt--;
          }

          if lex[o] == "(" {
            pCnt++;
          }
          if lex[o] == ")" {
            pCnt--;
          }

          if cbCnt != 0 || glCnt != 0 || bCnt != 0 || pCnt != 0 {
            exp_ = append(exp_, lex[o])
            continue
          }

          if lex[o] == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_, false)

        actions = append(actions, Action{ "print", "", exp_, exp, []string{}, []Action{}, []Condition{}, 6, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(2 + len(exp_))
      case "{":
        exp_ := []string{}

        //getting nb semicolons
        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i; o < len_lex; o++ {

          if lex[o] == "{" {
            cbCnt++;
          }
          if lex[o] == "}" {
            cbCnt--;
          }

          if lex[o] == "[:" {
            glCnt++;
          }
          if lex[o] == ":]" {
            glCnt--;
          }

          if lex[o] == "[" {
            bCnt++;
          }
          if lex[o] == "]" {
            bCnt--;
          }

          if lex[o] == "(" {
            pCnt++;
          }
          if lex[o] == ")" {
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

          if lex[o] == "newlineS" {
            break
          }
        }

        exp_ = exp_[1:]
        exp_ = exp_[:len(exp_) - 1]

        exp := actionizer(exp_, false)

        actions = append(actions, Action{ "group", "", []string{}, exp, []string{}, []Action{}, []Condition{}, 9, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(len(exp_) + 1)
      case "process":

        putFalsey := make(map[string][]Action)
        putFalsey["falsey"] = []Action{ Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString("undefined"))) }, false), false } }

        if lex[i + 1] == "~" {
          var procName = lex[i + 2]
          params := []string{}

          for o := i + 4; o < len_lex; o++ {
            if lex[o] == ")" {
              break
            }

            if lex[o] == "," {
              continue
            }

            params = append(params, lex[o])
          }
          i+=(len(params) + 5)

          var logic_ = []string{}

          cbCnt := 0

          for o := i; o < len_lex; o++ {
            if lex[o] == "{" {
              cbCnt++
            }

            if lex[o] == "}" {
              cbCnt--
            }

            logic_ = append(logic_, lex[o])

            if cbCnt == 0 {
              break
            }
          }

          i+=len(logic_) - 1

          logic := actionizer(logic_, false)

          actions = append(actions, Action{ "process", procName, []string{}, logic, params, []Action{}, []Condition{}, 10, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, putFalsey, actionizer([]string{ "�process" }, false), false })
        } else {
          params := []string{}

          for o := i + 2; o < len_lex; o+=2 {
            if lex[o] == ")" {
              break
            }

            params = append(params, lex[o])
          }
          i+=(3 + len(params))

          var logic_ = []string{}

          cbCnt := 0

          for o := i; o < len_lex; o++ {
            if lex[o] == "{" {
              cbCnt++
            }

            if lex[o] == "}" {
              cbCnt--
            }

            logic_ = append(logic_, lex[o])

            if cbCnt == 0 {
              break
            }
          }

          i+=len(logic_) - 1

          logic := actionizer(logic_, false)

          actions = append(actions, Action{ "process", "", []string{}, logic, params, []Action{}, []Condition{}, 10, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, putFalsey, actionizer([]string{ "�process" }, false), false })
        }
      case "#":

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 1

        var name = lex[i + 2]

        indexes := [][]string{[]string{}}
        var putIndexes [][]Action

        if lex[i + 3] == "." {

          cbCnt = 0
          glCnt = 0
          bCnt = 0
          pCnt = 0

          for o := i + 4; o < len_lex; o++ {
            if lex[o] == "{" {
              cbCnt++
            }
            if lex[o] == "[:" {
              glCnt++
            }
            if lex[o] == "[" {
              bCnt++
            }
            if lex[o] == "(" {
              pCnt++
            }

            if lex[o] == "}" {
              cbCnt--
            }
            if lex[o] == ":]" {
              glCnt--
            }
            if lex[o] == "]" {
              bCnt--
            }
            if lex[o] == ")" {
              pCnt--
            }

            if lex[o] == "." {
              indexes = append(indexes, []string{})
            } else {

              i++

              indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                if o < len_lex - 1 && lex[o + 1] == "." {
                  continue
                } else {
                  break
                }

              }
            }
          }

          for _, v := range indexes {

            v = v[1:len(v) - 1]
            putIndexes = append(putIndexes, actionizer(v, true))
          }

          i+=3
        }

        params := [][]string{{}}

        for o := i + 3; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++;
          }
          if lex[o] == "}" {
            cbCnt--;
          }

          if lex[o] == "[:" {
            glCnt++;
          }
          if lex[o] == ":]" {
            glCnt--;
          }

          if lex[o] == "[" {
            bCnt++;
          }
          if lex[o] == "]" {
            bCnt--;
          }

          if lex[o] == "(" {
            pCnt++;
          }
          if lex[o] == ")" {
            pCnt--;
          }

          if lex[o] == "(" {
            continue
          }

          if cbCnt != 0 && glCnt != 0 && bCnt != 0 && pCnt != 0 {
            params = append(params, []string{})
            continue
          }

          if lex[o] == ")" {
            break
          }

          if lex[o] == "," {
            params = append(params, []string{})
            continue
          }

          params[len(params) - 1] = append(params[len(params) - 1], lex[o])
        }

        var params_ = []Action{}

        for o := 0; o < len(params); o++ {
          params_ = append(params_, actionizer(params[o], true)...)
        }

        pCnt_ := 0
        skip_nums := 0

        for o := i; o < len_lex; o++ {
          if lex[o] == "(" {
            pCnt_++
          }
          if lex[o] == ")" {
            pCnt_--
          }

          skip_nums++;

          if pCnt_ == 0 && lex[o] == "newlineS" {
            break
          }
        }

        actions = append(actions, Action{ "#", name, []string{}, []Action{}, []string{}, params_, []Condition{}, 11, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=skip_nums
      case "return":

        returner_ := []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
            break
          }

          returner_ = append(returner_, lex[o])
        }

        returner := actionizer(returner_, true)

        actions = append(actions, Action{ "return", "", []string{}, returner, []string{}, []Action{}, []Condition{}, 12, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=len(returner_) + 2
      case "if":

        conditions := []Condition{}

        for ;lex[i] == "if"; {
          var cond_ = []string{}
          pCnt := 0

          for o := i + 1; o < len_lex; o++ {
            if lex[o] == "(" {
              pCnt++
            }
            if lex[o] == ")" {
              pCnt--
            }

            cond_ = append(cond_, lex[o])

            if pCnt == 0 {
              break
            }
          }

          cond := actionizer(cond_, true)
          cbCnt := 0
          actions_ := []string{}

          for o := i + 1 + len(cond_); o < len_lex; o++ {
            if lex[o] == "{" {
              cbCnt++
            }
            if lex[o] == "}" {
              cbCnt--
            }

            actions_ = append(actions_, lex[o])

            if cbCnt == 0 {
              break
            }
          }

          actions := actionizer(actions_, false)

          var if_ = Condition{ "if", cond, actions }

          conditions = append(conditions, if_)

          i+=(1 + len(cond_) + len(actions_))

          if i >= len_lex {
            break
          }
        }

        if !(i >= len_lex) {
          for ;lex[i] == "elseif"; {

            var cond_ = []string{}
            pCnt := 0

            for o := i + 1; o < len_lex; o++ {
              if lex[o] == "(" {
                pCnt++
              }
              if lex[o] == ")" {
                pCnt--
              }

              cond_ = append(cond_, lex[o])

              if pCnt == 0 {
                break
              }
            }

            cond := actionizer(cond_, true)
            cbCnt := 0
            actions_ := []string{}

            for o := i + 1 + len(cond_); o < len_lex; o++ {
              if lex[o] == "{" {
                cbCnt++
              }
              if lex[o] == "}" {
                cbCnt--
              }

              actions_ = append(actions_, lex[o])

              if cbCnt == 0 {
                break
              }
            }

            actions := actionizer(actions_, false)

            var elseif_ = Condition{ "elseif", cond, actions }

            conditions = append(conditions, elseif_)

            i+=(1 + len(cond_) + len(actions_))

            if i >= len_lex {
              break
            }
          }
        }

        if !(i >= len_lex) {
          actions_ := []string{}

          for ;lex[i] == "else"; {
            cbCnt := 0

            for o := i + 1; o < len_lex; o++ {
              if lex[o] == "{" {
                cbCnt++
              }
              if lex[o] == "}" {
                cbCnt--
              }

              if cbCnt == 0 {
                break
              }

              actions_ = append(actions_, lex[o])
            }

            actions := actionizer(actions_, false)

            var else_ = Condition{ "else", []Action{}, actions }

            conditions = append(conditions, else_)

            i+=(1 + len(actions_))

            if i >= len_lex {
              break
            }
          }
        }

        actions = append(actions, Action{ "conditional", "", []string{}, []Action{}, []string{}, []Action{}, conditions, 13, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
      case "import":

        var file = []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
            break
          }

          file = append(file, lex[o])
        }

        actionizedFile := actionizer(file, false)
        actions = append(actions, Action{ "import", "", []string{}, actionizedFile, []string{}, []Action{}, []Condition{}, 14, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(2 + len(file))
      case "read":
        var phrase = []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionizedPhrase := actionizer(phrase, true)
        actions = append(actions, Action{ "read", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 15, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(2 + len(phrase))
      case "break":
        actions = append(actions, Action{ "break", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 16, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
      case "skip":
        actions = append(actions, Action{ "skip", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 17, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
      case "eval":
        var phrase = []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionized := actionizer(phrase, true)
        actions = append(actions, Action{ "eval", "", []string{}, actionized, []string{}, []Action{}, []Condition{}, 18, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(2 + len(phrase))
      case "typeof":
        var phrase = []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionizedPhrase := actionizer(phrase, true)
        actions = append(actions, Action{ "typeof", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 19, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(2 + len(phrase))
      case "err":
        var phrase = []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionizedPhrase := actionizer(phrase, true)
        actions = append(actions, Action{ "err", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 20, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(2 + len(phrase))
      case "loop":

        var condition_ = []string{}

        pCnt := 0

        for o := i + 1; o < len_lex; o++ {
          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          condition_ = append(condition_, lex[o])

          if pCnt == 0 {
            break
          }
        }

        condition := actionizer(condition_, true)
        action_ := []string{}

        cbCnt := 0

        for o := i + 1 + len(condition_); o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          action_ = append(action_, lex[o])

          if cbCnt == 0 {
            break
          }
        }

        action := actionizer(action_, false)

        actions = append(actions, Action{ "loop", "", []string{}, action, []string{}, []Action{}, []Condition{ { "loop", condition, action } }, 21, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(1 + len(condition_) + len(action_))
      case "[:":
        var phrase = []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          phrase = append(phrase, lex[o])

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
            break
          }
        }

        i+=len(phrase)

        phrase = phrase[1:len(phrase) - 1]

        var _translated =  [][][]string{ [][]string{ []string{}, []string{} } }

        cbCnt = 0
        glCnt = 0
        bCnt = 0
        pCnt = 0

        cur := 0

        for _, k := range phrase {

          if k == "newlineN" {
            continue
          }

          if k == "{" {
            cbCnt++
          }
          if k == "}" {
            cbCnt--
          }

          if k == "[:" {
            glCnt++
          }
          if k == ":]" {
            glCnt--
          }

          if k == "[" {
            bCnt++
          }
          if k == "]" {
            bCnt--
          }

          if k == "(" {
            pCnt++
          }
          if k == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && k == ":" {
            cur = 1
            continue
          }
          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && k == "," {
            cur = 0
            _translated = append(_translated, [][]string{ []string{}, []string{} })
            continue
          }

          _translated[len(_translated) - 1][cur] = append(_translated[len(_translated) - 1][cur], k)
        }

        var translated = make(map[string][]Action)

        for _, k := range _translated {

          if len(k[0]) <= 0 {
            break
          }

          translated[k[0][0]] = actionizer(k[1], true)
        }

        if i >= len_lex {
          actions = append(actions, Action{ "hash", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, translated, actionizer([]string{ "�hash" }, false), false })
          break
        }

        isMutable := false

        //checks for a runtime hash
        for ;i < len_lex && lex[i] == ":::"; i++ {
          isMutable = !isMutable
        }

        if i >= len_lex {
          actions = append(actions, Action{ "hash", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, translated, actionizer([]string{ "�hash" }, false), isMutable })
          break
        }

        if lex[i] == "." {
          indexes := [][]string{[]string{}}

          cbCnt = 0
          glCnt = 0
          bCnt = 0
          pCnt = 0

          for o := i + 1; o < len_lex; o++ {
            if lex[o] == "{" {
              cbCnt++
            }
            if lex[o] == "[:" {
              glCnt++
            }
            if lex[o] == "[" {
              bCnt++
            }
            if lex[o] == "(" {
              pCnt++
            }

            if lex[o] == "}" {
              cbCnt--
            }
            if lex[o] == ":]" {
              glCnt--
            }
            if lex[o] == "]" {
              bCnt--
            }
            if lex[o] == ")" {
              pCnt--
            }

            if lex[o] == "." {
              indexes = append(indexes, []string{})
            } else {

              i++

              indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                if o < len_lex - 1 && lex[o + 1] == "." {
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
            putIndexes = append(putIndexes, actionizer(v, true))
          }

          i+=3

          actions = append(actions, Action{ "hashIndex", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 23, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, translated, actionizer([]string{ "�hash" }, false), isMutable })
        } else {
          actions = append(actions, Action{ "hash", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, translated, actionizer([]string{ "�hash" }, false), isMutable })
        }
      case "[":
        var phrase = []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
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

          var sub []string

          cbCnt := 0
          glCnt := 0
          bCnt := 0
          pCnt := 0

          for j := o; j < len(phrase); j++ {

            if phrase[j] == "{" {
              cbCnt++
            }
            if phrase[j] == "}" {
              cbCnt--
            }

            if phrase[j] == "[:" {
              glCnt++
            }
            if phrase[j] == ":]" {
              glCnt--
            }

            if phrase[j] == "[" {
              bCnt++
            }
            if phrase[j] == "]" {
              bCnt--
            }

            if phrase[j] == "(" {
              pCnt++
            }
            if phrase[j] == ")" {
              pCnt--
            }

            if phrase[j] == "," && cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {
              break
            }
            sub = append(sub, phrase[j])
          }

          o+=len(sub)

          arr = append(arr, actionizer(sub, true))
        }

        hashedArr := make(map[string][]Action)

        cur := "0"

        for _, v := range arr {
          hashedArr[cur] = v
          cur = add(cur, "1", paramCalcOpts{}, -1)
        }

        if i >= len_lex {
          actions = append(actions, Action{ "array", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, hashedArr, actionizer([]string{ "�array" }, false), false })
          break
        }

        isMutable := false

        //checks for a runtime array
        for ;i < len_lex && lex[i] == ":::"; i++ {
          isMutable = !isMutable
        }

        if i >= len_lex {
          actions = append(actions, Action{ "array", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, hashedArr, actionizer([]string{ "�array" }, false), isMutable })
          break
        }

        if lex[i] == "." {
          indexes := [][]string{[]string{}}

          cbCnt = 0
          glCnt = 0
          bCnt = 0
          pCnt = 0

          for o := i + 1; o < len_lex; o++ {
            if lex[o] == "{" {
              cbCnt++
            }
            if lex[o] == "[:" {
              glCnt++
            }
            if lex[o] == "[" {
              bCnt++
            }
            if lex[o] == "(" {
              pCnt++
            }

            if lex[o] == "}" {
              cbCnt--
            }
            if lex[o] == ":]" {
              glCnt--
            }
            if lex[o] == "]" {
              bCnt--
            }
            if lex[o] == ")" {
              pCnt--
            }

            if lex[o] == "." {
              indexes = append(indexes, []string{})
            } else {

              i++

              indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                if o < len_lex - 1 && lex[o + 1] == "." {
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
            putIndexes = append(putIndexes, actionizer(v, true))
          }

          i+=3

          actions = append(actions, Action{ "arrayIndex", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 25, []Action{}, []Action{}, []Action{}, arr, putIndexes, hashedArr, actionizer([]string{ "�array" }, false), isMutable })
        } else {
          actions = append(actions, Action{ "array", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, hashedArr, actionizer([]string{ "�array" }, false), isMutable })
        }
      case "ascii":
        var phrase = []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionizedPhrase := actionizer(phrase, true)
        actions = append(actions, Action{ "ascii", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 26, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(2 + len(phrase))
      case "parse":
        var phrase = []string{}

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 0

        for o := i + 2; o < len_lex; o++ {
          if lex[o] == "{" {
            cbCnt++
          }
          if lex[o] == "}" {
            cbCnt--
          }

          if lex[o] == "[:" {
            glCnt++
          }
          if lex[o] == ":]" {
            glCnt--
          }

          if lex[o] == "[" {
            bCnt++
          }
          if lex[o] == "]" {
            bCnt--
          }

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
            break
          }

          phrase = append(phrase, lex[o])
        }

        actionizedPhrase := actionizer(phrase, true)
        actions = append(actions, Action{ "parse", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 27, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=(2 + len(phrase))
      case "len":

        var exp []string
        pCnt := 0

        for o := i + 1; o < len_lex; o++ {
          if lex[o] == "(" {
            pCnt++
            continue
          }
          if lex[o] == ")" {
            pCnt--
            continue
          }

          if pCnt == 0 {
            break
          }

          exp = append(exp, lex[o])
        }

        actionized := actionizer(exp, true)

        actions = append(actions, Action{ "len", "", []string{}, actionized, []string{}, []Action{}, []Condition{}, 31, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
        i+=3 + len(exp)
      default:

        valPuts := func(lex []string, i int) int {

          //KEEP IN MIND: type key starts with ascii of 233
          //KEEP IN MIND: index key starts with ascii of 8

          if i >= len_lex {
            return 1
          }

          isMutable := false
          val := lex[i]

          //checks for a runtime value
          for ;i + 1 < len_lex && lex[i + 1] == ":::"; i++ {
            isMutable = !isMutable
          }

          i++

          switch C.GoString(GetType(C.CString(val))) {

            case "string": {

              noQ := []rune(val)[1:len(val) - 1]
              hashedString := make(map[string][]Action)

              //specify the value for the "falsey" case
              hashedString["falsey"] = []Action{ Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString("undefined"))) }, false), isMutable  } }

              cur := "0"

              for _, v := range noQ {
                hashedString[cur] = actionizer([]string{ "" + string(v) }, true)
                cur = add(cur, "1", paramCalcOpts{}, -1)
              }

              actions = append(actions, Action{ "string", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashedString, actionizer([]string{ "�" + C.GoString(GetType(C.CString(val))) }, false), isMutable })
            }
            case "number":

              hashed := make(map[string][]Action)

              //specify the value for the "falsey" case
              hashed["falsey"] = []Action{ Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString("undefined"))) }, false), isMutable } }

              actions = append(actions, Action{ "number", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, actionizer([]string{ "�" + C.GoString(GetType(C.CString(val))) }, false), false })
            case "boolean":

              hashed := make(map[string][]Action)

              //specify the value for the "falsey" case
              hashed["falsey"] = []Action{ Action{ "boolean", "", []string{ strconv.FormatBool(val != "true") }, []Action{}, []string{}, []Action{}, []Condition{}, 40, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString("boolean"))) }, false), isMutable } }

              actions = append(actions, Action{ "boolean", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 40, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, actionizer([]string{ "�" + C.GoString(GetType(C.CString(val))) }, false), isMutable })
            case "falsey":

              hashed := make(map[string][]Action)

              //specify the value for the "falsey" case
              hashed["falsey"] = []Action{ Action{ "falsey", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString("undefined"))) }, false), isMutable } }

              actions = append(actions, Action{ "falsey", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, actionizer([]string{ "�" + C.GoString(GetType(C.CString(val))) }, false), isMutable })
            case "none":

              if strings.HasPrefix(val, "$") {

                actions = append(actions, Action{ "variable", val, []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 43, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�variable" }, false), isMutable })
              } else if strings.HasPrefix(val, "�") {

                actions = append(actions, Action{ "type", "", []string{ strings.TrimPrefix(val, "�") }, []Action{}, []string{}, []Action{}, []Condition{}, GetActNum(strings.TrimPrefix(val, "�")), []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), []Action{}, false })
              } else if strings.HasPrefix(val, "") {

                hashed := make(map[string][]Action)
                hashed["falsey"] = []Action{ Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString("undefined"))) }, false), false } }

                actions = append(actions, Action{ "index_key", "", []string{ strings.TrimPrefix(val, "") }, []Action{}, []string{}, []Action{}, []Condition{}, GetActNum(strings.TrimPrefix(val, "")), []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, hashed, []Action{}, false })
              } else {

                hashedString := make(map[string][]Action)

                //specify the value for the "falsey" case
                hashedString["falsey"] = []Action{ Action{ "falsey", "", []string{ "undefined" }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString("undefined"))) }, false), isMutable } }

                //get it? 42?
                actions = append(actions, Action{ "none", "", []string{ val }, []Action{}, []string{}, []Action{}, []Condition{}, 42, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString(val))) }, false), isMutable })
              }
          }

          return 0
        }

        if i + 1 < len_lex {

          if (lex[i + 1] == "++" || lex[i + 1] == "--") && strings.HasPrefix(lex[i], "$") {

            id_ := []byte(lex[i + 1])

            id := ""

            for o := 0; o < len(id_); o++ {
              _id := strconv.Itoa(int(id_[o]))
              id+=_id
            }

            intID, _ := strconv.Atoi(id)

            actions = append(actions, Action{ lex[i + 1], lex[i], []string{}, []Action{}, []string{}, []Action{}, []Condition{}, intID, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
            i++
            continue;
          }

          if (lex[i + 1] == "+=" || lex[i + 1] == "-=" || lex[i + 1] == "*=" || lex[i + 1] == "/=" || lex[i + 1] == "%=" || lex[i + 1] == "^=") && strings.HasPrefix(lex[i], "$") {

            var by_ []string

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 2; o < len_lex; o++ {
              if lex[o] == "{" {
                cbCnt++
              }
              if lex[o] == "}" {
                cbCnt--
              }

              if lex[o] == "[:" {
                glCnt++
              }
              if lex[o] == ":]" {
                glCnt--
              }

              if lex[o] == "[" {
                bCnt++
              }
              if lex[o] == "]" {
                bCnt--
              }

              if lex[o] == "(" {
                pCnt++
              }
              if lex[o] == ")" {
                pCnt--
              }

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
                break
              }

              by_ = append(by_, lex[o])
            }

            by := actionizer(by_, true)

            id_ := []byte(lex[i + 1])

            id := ""

            for o := 0; o < len(id_); o++ {
              _id := strconv.Itoa(int(id_[o]))
              id+=_id
            }

            intID, _ := strconv.Atoi(id)

            actions = append(actions, Action{ lex[i + 1], lex[i], []string{}, by, []string{}, []Action{}, []Condition{}, intID, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
            continue;
          }

          doPutIndex := false

          icbCnt := 0
          iglCnt := 0
          ibCnt := 0
          ipCnt := 0

          for o := i; o < len_lex; o++ {
            if lex[o] == "{" {
              icbCnt++
            }
            if lex[o] == "}" {
              icbCnt--
            }

            if lex[o] == "[:" {
              iglCnt++
            }
            if lex[o] == ":]" {
              iglCnt--
            }

            if lex[o] == "[" {
              ibCnt++
            }
            if lex[o] == "]" {
              ibCnt--
            }

            if lex[o] == "(" {
              ipCnt++
            }
            if lex[o] == ")" {
              ipCnt--
            }

            if icbCnt == 0 && iglCnt == 0 && ibCnt == 0 && ipCnt == 0 && lex[o] == "newlineS" {
              break
            }

            if icbCnt == 0 && iglCnt == 0 && ibCnt == 0 && ipCnt == 0 && lex[o] == ":" {
              doPutIndex = true
              break
            }
          }

          var indexes [][]Action
          name := lex[i]

          if lex[i + 1] == "." && lex[i + 2] == "[" && doPutIndex {

            _indexes := [][]string{}

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 1; o < len_lex; i, o = i + 1, o + 1 {
              if lex[o] == "{" {
                cbCnt++
              }
              if lex[o] == "}" {
                cbCnt--
              }

              if lex[o] == "[:" {
                glCnt++
              }
              if lex[o] == ":]" {
                glCnt--
              }

              if lex[o] == "[" {
                bCnt++
              }
              if lex[o] == "]" {
                bCnt--
              }

              if lex[o] == "(" {
                pCnt++
              }
              if lex[o] == ")" {
                pCnt--
              }

              if lex[o] == "." {
                _indexes = append(_indexes, []string{})
                continue
              }

              _indexes[len(_indexes) - 1] = append(_indexes[len(_indexes) - 1], lex[o])

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o + 1] == ":" {
                break
              }
            }

            for _, v := range _indexes {
              indexes = append(indexes, actionizer(v[1:len(v) - 1], true))
            }

            i++
          }

          if lex[i + 1] == ":" && (strings.HasPrefix(lex[i], "$") || lex[i] == "]") {
            exp_ := []string{}

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            for o := i + 2; o < len_lex; o++ {

              if lex[o] == "{" {
                cbCnt++
              }
              if lex[o] == "}" {
                cbCnt--
              }

              if lex[o] == "[:" {
                glCnt++
              }
              if lex[o] == ":]" {
                glCnt--
              }

              if lex[o] == "[" {
                bCnt++
              }
              if lex[o] == "]" {
                bCnt--
              }

              if lex[o] == "(" {
                pCnt++
              }
              if lex[o] == ")" {
                pCnt--
              }

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
                break
              }

              exp_ = append(exp_, lex[o]);
            }

            exp := actionizer(exp_, false)

            actions = append(actions, Action{ "let", name, exp_, exp, []string{}, []Action{}, []Condition{}, 28, []Action{}, []Action{}, []Action{}, [][]Action{}, indexes, make(map[string][]Action), actionizer([]string{ "�statement" }, false), false })
            i+=(len(exp))
            continue
          }

          if lex[i + 1] == "." {

            val := lex[i]

            cbCnt := 0
            glCnt := 0
            bCnt := 0
            pCnt := 0

            indexes := [][]string{[]string{}}

            cbCnt = 0
            glCnt = 0
            bCnt = 0
            pCnt = 0

            for o := i + 2; o < len_lex; o++ {
              if lex[o] == "{" {
                cbCnt++
              }
              if lex[o] == "[:" {
                glCnt++
              }
              if lex[o] == "[" {
                bCnt++
              }
              if lex[o] == "(" {
                pCnt++
              }

              if lex[o] == "}" {
                cbCnt--
              }
              if lex[o] == ":]" {
                glCnt--
              }
              if lex[o] == "]" {
                bCnt--
              }
              if lex[o] == ")" {
                pCnt--
              }

              if lex[o] == "." {
                indexes = append(indexes, []string{})
              } else {

                i++

                indexes[len(indexes) - 1] = append(indexes[len(indexes) - 1], lex[o])

                if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 {

                  if o < len_lex - 1 && lex[o + 1] == "." {
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
              putIndexes = append(putIndexes, actionizer(v, true))
            }

            i+=3

            if strings.HasPrefix(val, "$") {
              actVal := actionizer([]string{ val }, true)

              actions = append(actions, Action{ "variableIndex", "", []string{}, actVal, []string{}, []Action{}, []Condition{}, 46, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString(val))) }, false), false })
            } else {
              actVal := actionizer([]string{ val }, true)

              actions = append(actions, Action{ "expressionIndex", "", []string{}, actVal, []string{}, []Action{}, []Condition{}, 8, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString(val))) }, false), false })
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
