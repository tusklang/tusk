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
  Index_Type      string
  Hash_Values     map[string][]Action

  //type of the action as a value
  ValueType     []Action
}

func actionizer(lex []string) []Action {
  var actions = []Action{}
  var len_lex = len(lex)

  var operations = []string{ "+", "-", "*", "/", "^", "%", "&", "|", "=", ">", "<", ">=", "<=", ")", "(", "~~", "~~~", "!" }

  for i := 0; i < len_lex; i++ {

    oCbCnt := 0
    oGlCnt := 0
    oBCnt := 0
    oPCnt := 0

    doExpress := false

    for o := i; o < len_lex; o++ {
      if lex[o] == "{" {
        oCbCnt++;
      }
      if lex[o] == "}" {
        oCbCnt--;
      }

      if lex[o] == "[:" {
        oGlCnt++;
      }
      if lex[o] == ":]" {
        oGlCnt--;
      }

      if lex[o] == "[" {
        oBCnt++;
      }
      if lex[o] == "]" {
        oBCnt--;
      }

      if lex[o] == "(" {
        oPCnt++;
      }
      if lex[o] == ")" {
        oPCnt--;
      }

      if lex[o] == "." {
        continue;
      }

      if oCbCnt != 0 && oGlCnt != 0 && oBCnt != 0 && oPCnt != 0 && (arrayContain([]string{ "+", "-", "*", "/", "^", "%", "=", ">", "<", ">=", "<=", "&", "|", "~~", "~~~" }, lex[o]) || arrayContain([]string{ "!", "(" }, lex[i])) {
        doExpress = true
        break
      }
    }

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

        pExpAct := actionizer(pExp)

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

          num1 = actionizer(num)

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

          num2 = actionizer(num)

        } else {

          for _, v := range _num2 {
            num2 = append(num2, v.(Action))
          }

        }

        var act_exp = Action{ "exponentiate", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 36, num1, num2, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�operation" }) }

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

            num1 = actionizer(num)

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

            num2 = actionizer(num)

          } else {

            for _, v := range _num2 {
              num2 = append(num2, v.(Action))
            }

          }

          var act_exp = Action{ "multiply", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 34, num1, num2, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�operation" }) }

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

            num1 = actionizer(num)

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

            num2 = actionizer(num)

          } else {

            for _, v := range _num2 {
              num2 = append(num2, v.(Action))
            }

          }

          var act_exp = Action{ "divide", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 35, num1, num2, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�operation" }) }

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

          num1 = actionizer(num)

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

          num2 = actionizer(num)

        } else {

          for _, v := range _num2 {
            num2 = append(num2, v.(Action))
          }

        }

        var act_exp = Action{ "modulo", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 37, num1, num2, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�operation" }) }

        exp_ := append(exp[:index - 1], act_exp)
        exp_ = append(exp_, exp[index + 1:])

        exp = exp_
      }

      for ;interfaceContain(exp, "+") || interfaceContain(exp, "-"); {

        if interfaceIndexOf("+", exp) > interfaceIndexOf("-", exp) || interfaceIndexOf("-", exp) == -1 {
          index := interfaceIndexOf("+", exp)

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

            num1 = actionizer(num)

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

            num2 = actionizer(num)

          } else {

            for _, v := range _num2 {
              num2 = append(num2, v.(Action))
            }

          }

          var act_exp = Action{ "add", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 32, num1, num2, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�operation" }) }

          exp_ := append(exp[:index - 1], act_exp)
          exp_ = append(exp_, exp[index + 1:])

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

            num1 = actionizer(num)

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

            num2 = actionizer(num)

          } else {

            for _, v := range _num2 {
              num2 = append(num2, v.(Action))
            }

          }

          var act_exp = Action{ "subtract", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 33, num1, num2, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�operation" }) }

          exp_ := append(exp[:index - 1], act_exp)
          exp_ = append(exp_, exp[index + 1:])

          exp = exp_
        }

      }

      actions = append(actions, exp[0].(Action))
    }

    if i >= len_lex {
      break
    }

    switch lex[i] {
      case "newlineN":
        actions = append(actions, Action{ "newline", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 0, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        exp := actionizer(exp_)

        actions = append(actions, Action{ "local", lex[i + 2], []string{}, exp, []string{}, []Action{}, []Condition{}, 1, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        exp := actionizer(exp_)

        actions = append(actions, Action{ "dynamic", lex[i + 2], []string{}, exp, []string{}, []Action{}, []Condition{}, 2, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
        i+=(4 + len(exp_))
      case "alt":

        var alter = Action{ "alt", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 3, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) }

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

        cond := actionizer(cond_)

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
          actions := actionizer(actions_)

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

        exp := actionizer(exp_)

        actions = append(actions, Action{ "global", lex[i + 2], []string{}, exp, []string{}, []Action{}, []Condition{}, 4, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        exp := actionizer(exp_)

        actions = append(actions, Action{ "log", "", exp_, exp, []string{}, []Action{}, []Condition{}, 5, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        exp := actionizer(exp_)

        actions = append(actions, Action{ "print", "", exp_, exp, []string{}, []Action{}, []Condition{}, 6, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        exp := actionizer(exp_)

        actions = append(actions, Action{ "group", "", []string{}, exp, []string{}, []Action{}, []Condition{}, 9, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
        i+=(len(exp_) + 1)
      case "process":
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

          logic := actionizer(logic_)

          actions = append(actions, Action{ "process", procName, []string{}, logic, params, []Action{}, []Condition{}, 10, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

          logic := actionizer(logic_)

          actions = append(actions, Action{ "process", "", []string{}, logic, params, []Action{}, []Condition{}, 10, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
        }
      case "#":

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 1

        var name = lex[i + 2]

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
          params_ = append(params_, actionizer(params[o])...)
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

        actions = append(actions, Action{ "#", name, []string{}, []Action{}, []string{}, params_, []Condition{}, 11, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        returner := actionizer(returner_)

        actions = append(actions, Action{ "return", "", []string{}, returner, []string{}, []Action{}, []Condition{}, 12, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
        i+=len(returner) + 2
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

          cond := actionizer(cond_)
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

          actions := actionizer(actions_)

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

            cond := actionizer(cond_)
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

            actions := actionizer(actions_)

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

            actions := actionizer(actions_)

            var else_ = Condition{ "else", []Action{}, actions }

            conditions = append(conditions, else_)

            i+=(1 + len(actions_))

            if i >= len_lex {
              break
            }
          }
        }

        actions = append(actions, Action{ "conditional", "", []string{}, []Action{}, []string{}, []Action{}, conditions, 13, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        actionizedFile := actionizer(file)
        actions = append(actions, Action{ "import", "", []string{}, actionizedFile, []string{}, []Action{}, []Condition{}, 14, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        actionizedPhrase := actionizer(phrase)
        actions = append(actions, Action{ "read", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 15, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
        i+=(2 + len(phrase))
      case "break":
        actions = append(actions, Action{ "break", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 16, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
      case "skip":
        actions = append(actions, Action{ "skip", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 17, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        actionized := actionizer(phrase)
        actions = append(actions, Action{ "eval", "", []string{}, actionized, []string{}, []Action{}, []Condition{}, 18, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        actionizedPhrase := actionizer(phrase)
        actions = append(actions, Action{ "typeof", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 19, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        actionizedPhrase := actionizer(phrase)
        actions = append(actions, Action{ "err", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 20, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        condition := actionizer(condition_)
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

        action := actionizer(action_)

        actions = append(actions, Action{ "loop", "", []string{}, action, []string{}, []Action{}, []Condition{ { "loop", condition, action } }, 21, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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
          translated[k[0][0]] = actionizer(k[1])
        }

        if i >= len_lex {
          actions = append(actions, Action{ "hash", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "hash", translated, actionizer([]string{ "�hash" }) })
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
            putIndexes = append(putIndexes, actionizer(v))
          }

          i+=3

          actions = append(actions, Action{ "hashIndex", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 23, []Action{}, []Action{}, []Action{}, [][]Action{}, putIndexes, "hash", translated, actionizer([]string{ "�hash" }) })
        } else {
          actions = append(actions, Action{ "hash", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 22, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "hash", translated, actionizer([]string{ "�hash" }) })
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

          arr = append(arr, actionizer(sub))
        }

        if i >= len_lex {
          actions = append(actions, Action{ "array", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, "array", make(map[string][]Action), actionizer([]string{ "�array" }) })
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
            putIndexes = append(putIndexes, actionizer(v))
          }

          i+=3

          actions = append(actions, Action{ "arrayIndex", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, 25, []Action{}, []Action{}, []Action{}, arr, putIndexes, "array", make(map[string][]Action), actionizer([]string{ "�array" }) })
        } else {
          actions = append(actions, Action{ "array", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, 24, []Action{}, []Action{}, []Action{}, arr, [][]Action{}, "array", make(map[string][]Action), actionizer([]string{ "�array" }) })
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

        actionizedPhrase := actionizer(phrase)
        actions = append(actions, Action{ "ascii", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 26, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        actionizedPhrase := actionizer(phrase)
        actions = append(actions, Action{ "parse", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, 27, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

        actionized := actionizer(exp)

        actions = append(actions, Action{ "len", "", []string{}, actionized, []string{}, []Action{}, []Condition{}, 31, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
        i+=3 + len(exp)
      default:
        if i + 1 < len_lex {

          if (lex[i + 1] == "++" || lex[i + 1] == "--") && strings.HasPrefix(lex[i], "$") {

            id_ := []byte(lex[i + 1])

            id := ""

            for o := 0; o < len(id_); o++ {
              _id := strconv.Itoa(int(id_[o]))
              id+=_id
            }

            intID, _ := strconv.Atoi(id)

            actions = append(actions, Action{ lex[i + 1], lex[i], []string{}, []Action{}, []string{}, []Action{}, []Condition{}, intID, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
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

            by := actionizer(by_)

            id_ := []byte(lex[i + 1])

            id := ""

            for o := 0; o < len(id_); o++ {
              _id := strconv.Itoa(int(id_[o]))
              id+=_id
            }

            intID, _ := strconv.Atoi(id)

            actions = append(actions, Action{ lex[i + 1], lex[i], []string{}, by, []string{}, []Action{}, []Condition{}, intID, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
            continue;
          }

          if lex[i + 1] == ":" && strings.HasPrefix(lex[i], "$") {
            exp_ := []string{}

            for o := i + 2; o < len_lex; o++ {

              if lex[o] == "newlineS" {
                break;
              }

              exp_ = append(exp_, lex[o]);
            }

            exp := actionizer(exp_)

            actions = append(actions, Action{ "let", lex[i], exp_, exp, []string{}, []Action{}, []Condition{}, 28, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�statement" }) })
            i+=(len(exp))
          }
        } else {

          if i + 1 < len_lex && lex[i + 1] == "." {
            valAct := actionizer([]string{ lex[i] })

            var _indexes [][]string

            bCnt := 0

            for o := i + 1; o < len_lex; o++ {

              if lex[o] == "[" {
                bCnt++
              }
              if lex[o] == "]" {
                bCnt--
              }

              if bCnt != 0 {
                _indexes[len(_indexes) - 1] = append(_indexes[len(_indexes) - 1], lex[o])
              }

              if lex[o] == "." {

                _indexes = append(_indexes, []string{})

                continue
              }

              break
            }

            _ = valAct
          } else {
            //KEEP IN MIND: type key starts with ascii of 233

            switch C.GoString(GetType(C.CString(lex[i]))) {
              case "string":
                actions = append(actions, Action{ "string", "", []string{ lex[i] }, []Action{}, []string{}, []Action{}, []Condition{}, 38, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString(lex[i]))) }) })
              case "number":
                actions = append(actions, Action{ "number", "", []string{ lex[i] }, []Action{}, []string{}, []Action{}, []Condition{}, 39, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString(lex[i]))) }) })
              case "boolean":
                actions = append(actions, Action{ "boolean", "", []string{ lex[i] }, []Action{}, []string{}, []Action{}, []Condition{}, 40, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString(lex[i]))) }) })
              case "falsey":
                actions = append(actions, Action{ "falsey", "", []string{ lex[i] }, []Action{}, []string{}, []Action{}, []Condition{}, 41, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString(lex[i]))) }) })
              case "none":

                if strings.HasPrefix(lex[i], "$") {

                  actions = append(actions, Action{ "variable", "", []string{ lex[i] }, []Action{}, []string{}, []Action{}, []Condition{}, 43, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�variable" }) })
                } else if strings.HasPrefix(lex[i], "�") {

                  actions = append(actions, Action{ strings.TrimPrefix(lex[i], "�"), "", []string{ strings.TrimPrefix(lex[i], "�") }, []Action{}, []string{}, []Action{}, []Condition{}, GetActNum(strings.TrimPrefix(lex[i], "�")), []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), []Action{} })
                } else {

                  //get it? 42?
                  actions = append(actions, Action{ "none", "", []string{ lex[i] }, []Action{}, []string{}, []Action{}, []Condition{}, 42, []Action{}, []Action{}, []Action{}, [][]Action{}, [][]Action{}, "", make(map[string][]Action), actionizer([]string{ "�" + C.GoString(GetType(C.CString(lex[i]))) }) })
                }
            }
          }
        }
      }
  }

  return actions
}
