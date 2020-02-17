package main

import "strings"
import "regexp"

func replaceFrom(orig []string, start, end int, with string) []string {
  return append(append(orig[:start], with), orig[end + 1:]...)
}

func mathParse(gexp *[][]string, functions []Funcs, line uint64, calc_params paramCalcOpts, vars map[string]Variable, dir string) [][]string {

  exp := *gexp

  var orig = [][][]rune{[][]rune{}}

  for i := 0; i < len(exp); i++ {
    for o := 0; o < len(exp[i]); o++ {
      orig[len(orig) - 1] = append(orig[len(orig) - 1], []rune(exp[i][o]))
    }
  }

  if len(exp) == 0 {
    return [][]string{[]string{"0"}}
  } else if len(exp[0]) == 0 {
    return exp
  } else if exp[0][0] == "true" || exp[0][0] == "false" || exp[0][0] == "undefined" || exp[0][0] == "null" {
    return exp
  } else {
    for i := 0; i < len(exp); i++ {

      for o := 0; o < len(exp[i]); o++ {
        if strings.HasPrefix(exp[i][o], "$") {

          if o - 1 != -1 {
            if exp[i][o - 1] == "~" {
              continue
            }
          }

          if vars[exp[i][o]].Value == nil {
            exp[i][o] = parser(vars[exp[i][o][1:]].ValueActs, calc_params, dir, line, functions, vars, false).Exp[0][0]
          } else {
            exp[i][o] = vars[exp[i][o]].Value[0]
          }
        }
      }
    }

    for ;arrayContain2Nest(exp, "(") && arrayContain2Nest(exp, ")"); {

      start := indexOf2Nest("(", exp)
      end := indexOf2Nest(")", exp)
      parenExp := exp[start[0]][start[1] + 1:end[1]]

      actions := actionizer(parenExp)

      evaled := parser(actions, calc_params, dir, line, functions, vars, false).Exp[0][0]

      exp[start[0]] = replaceFrom(exp[start[0]], start[1], end[1], evaled)
    }

    for ;arrayContain2Nest(exp, "^"); {
      op := indexOf2Nest("^", exp)

      _o := exp[op[0]][op[1] - 1:op[1]]
      o_ := exp[op[0]][op[1] + 1:op[1] + 2]

      if len(exp[op[0]][:op[1]]) >= 2 && exp[op[0]][op[1] - 2:op[1]][0] == "-" {
        _o = exp[op[0]][op[1] - 2:op[1]]
      }

      if len(exp[op[0]][op[1] + 1:]) >= 2 && exp[op[0]][op[1] + 1:][0] == "-" {
        o_ = exp[op[0]][op[1] + 1:op[1] + 3]
      }

      var num = exponentiate(strings.Join(_o, ""), strings.Join(o_, ""), calc_params, line, functions)

      if len(_o) == 2 && len(o_) == 2 {
        exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 2, num)
      } else if len(_o) != 2 && len(o_) == 2 {
        exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 2, num)
      } else if len(_o) == 2 && len(o_) != 2 {
        exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 1, num)
      } else if len(_o) != 2 && len(o_) != 2 {
        exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 1, num)
      }
    }

    for ;arrayContain2Nest(exp, "*") || arrayContain2Nest(exp, "/"); {

      if (indexOf2Nest("*", exp)[0] <= indexOf2Nest("/", exp)[0] || indexOf2Nest("/", exp)[0] == -1) && (indexOf2Nest("*", exp)[0] != -1) {
        op := indexOf2Nest("*", exp)

        _o := exp[op[0]][op[1] - 1:op[1]]
        o_ := exp[op[0]][op[1] + 1:op[1] + 2]

        if len(exp[op[0]][:op[1]]) >= 2 && exp[op[0]][op[1] - 2:op[1]][0] == "-" {
          _o = exp[op[0]][op[1] - 2:op[1]]
        }

        if len(exp[op[0]][op[1] + 1:]) >= 2 && exp[op[0]][op[1] + 1:][0] == "-" {
          o_ = exp[op[0]][op[1] + 1:op[1] + 3]
        }

        var num = multiply(strings.Join(_o, ""), strings.Join(o_, ""), calc_params, line, functions)

        if len(_o) == 2 && len(o_) == 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 2, num)
        } else if len(_o) != 2 && len(o_) == 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 2, num)
        } else if len(_o) == 2 && len(o_) != 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 1, num)
        } else if len(_o) != 2 && len(o_) != 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 1, num)
        }
      } else {
        op := indexOf2Nest("/", exp)

        _o := exp[op[0]][op[1] - 1:op[1]]
        o_ := exp[op[0]][op[1] + 1:op[1] + 2]

        if len(exp[op[0]][:op[1]]) >= 2 && exp[op[0]][op[1] - 2:op[1]][0] == "-" {
          _o = exp[op[0]][op[1] - 2:op[1]]
        }

        if len(exp[op[0]][op[1] + 1:]) >= 2 && exp[op[0]][op[1] + 1:][0] == "-" {
          o_ = exp[op[0]][op[1] + 1:op[1] + 3]
        }

        var num = division(strings.Join(_o, ""), strings.Join(o_, ""), calc_params, line, functions)

        if len(_o) == 2 && len(o_) == 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 2, num)
        } else if len(_o) != 2 && len(o_) == 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 2, num)
        } else if len(_o) == 2 && len(o_) != 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 1, num)
        } else if len(_o) != 2 && len(o_) != 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 1, num)
        }
      }

    }

    for ;arrayContain2Nest(exp, "%"); {
      op := indexOf2Nest("%", exp)

      _o := exp[op[0]][op[1] - 1:op[1]]
      o_ := exp[op[0]][op[1] + 1:op[1] + 2]

      if len(exp[op[0]][:op[1]]) >= 2 && exp[op[0]][op[1] - 2:op[1]][0] == "-" {
        _o = exp[op[0]][op[1] - 2:op[1]]
      }

      if len(exp[op[0]][op[1] + 1:]) >= 2 && exp[op[0]][op[1] + 1:][0] == "-" {
        o_ = exp[op[0]][op[1] + 1:op[1] + 3]
      }

      var num = modulo(strings.Join(_o, ""), strings.Join(o_, ""), calc_params, line, functions)

      if len(_o) == 2 && len(o_) == 2 {
        exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 2, num)
      } else if len(_o) != 2 && len(o_) == 2 {
        exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 2, num)
      } else if len(_o) == 2 && len(o_) != 2 {
        exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 1, num)
      } else if len(_o) != 2 && len(o_) != 2 {
        exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 1, num)
      }
    }

    for ;arrayContain2Nest(exp, "+") || arrayContain2Nest(exp, "-"); {

      if (indexOf2Nest("+", exp)[0] <= indexOf2Nest("-", exp)[0] || indexOf2Nest("-", exp)[0] == -1) && indexOf2Nest("+", exp)[0] != -1 {
        op := indexOf2Nest("+", exp)

        _o := exp[op[0]][op[1] - 1:op[1]]
        o_ := exp[op[0]][op[1] + 1:op[1] + 2]

        if len(exp[op[0]][:op[1]]) >= 2 && exp[op[0]][op[1] - 2:op[1]][0] == "-" {
          _o = exp[op[0]][op[1] - 2:op[1]]
        }

        if len(exp[op[0]][op[1] + 1:]) >= 2 && exp[op[0]][op[1] + 1:][0] == "-" {
          o_ = exp[op[0]][op[1] + 1:op[1] + 3]
        }

        var num = add(strings.Join(_o, ""), strings.Join(o_, ""), calc_params, line, functions)

        if len(_o) == 2 && len(o_) == 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 2, num)
        } else if len(_o) != 2 && len(o_) == 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 2, num)
        } else if len(_o) == 2 && len(o_) != 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 1, num)
        } else if len(_o) != 2 && len(o_) != 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 1, num)
        }
      } else {
        op := indexOf2Nest("-", exp)

        last := ""

        for i := 0; i < len(exp[op[0]]); i++ {

          lastNum, _ := regexp.MatchString("\\d+", last)

          if lastNum {
            op[1] = i
          }

          last = exp[op[0]][i]
        }

        _o := exp[op[0]][op[1] - 1:op[1]]
        o_ := exp[op[0]][op[1] + 1:op[1] + 2]

        if len(exp[op[0]][:op[1]]) >= 2 && exp[op[0]][op[1] - 2:op[1]][0] == "-" {
          _o = exp[op[0]][op[1] - 2:op[1]]
        }

        if len(exp[op[0]][op[1] + 1:]) >= 2 && exp[op[0]][op[1] + 1:][0] == "-" {
          o_ = exp[op[0]][op[1] + 1:op[1] + 3]
        }

        var num = subtract(strings.Join(_o, ""), strings.Join(o_, ""), calc_params, line, functions)

        if len(_o) == 2 && len(o_) == 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 2, num)
        } else if len(_o) != 2 && len(o_) == 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 2, num)
        } else if len(_o) == 2 && len(o_) != 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 2, op[1] + 1, num)
        } else if len(_o) != 2 && len(o_) != 2 {
          exp[op[0]] = replaceFrom(exp[op[0]], op[1] - 1, op[1] + 1, num)
        }
      }

    }

    for ;arrayContain2Nest(exp, ">=") || arrayContain2Nest(exp, "<="); {
      if (indexOf2Nest(">=", exp)[0] < indexOf2Nest("<=", exp)[0] || indexOf2Nest("<=", exp)[0] == -1)  && indexOf2Nest(">=", exp)[0] != -1 {
        index := indexOf2Nest(">=", exp)

        if !isLess(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1]) {
          exp[index[0]] = []string{ "true" }
        } else {
          exp[index[0]] = []string{ "false" }
        }
      } else {
        index := indexOf2Nest("<=", exp)

        if isLess(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1]) || returnInit(exp[index[0]][index[1] - 1]) == returnInit(exp[index[0]][index[1] + 1]) {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "false")
        }
      }
    }

    for ;arrayContain2Nest(exp, ">") || arrayContain2Nest(exp, "<"); {
      if (indexOf2Nest(">", exp)[0] < indexOf2Nest("<", exp)[0] || indexOf2Nest("<", exp)[0] == -1) && indexOf2Nest(">", exp)[0] != -1 {
        index := indexOf2Nest(">", exp)

        if !isLess(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1]) && returnInit(exp[index[0]][index[1] - 1]) != returnInit(exp[index[0]][index[1] + 1]) {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "false")
        }
      } else {
        index := indexOf2Nest("<", exp)

        if isLess(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1]) {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "false")
        }
      }
    }

    for ;arrayContain2Nest(exp, "=") || arrayContain2Nest(exp, "!="); {
      if (indexOf2Nest("=", exp)[0] < indexOf2Nest("!=", exp)[0] || indexOf2Nest("!=", exp)[0] == -1) && indexOf2Nest("=", exp)[0] != -1 {
        index := indexOf2Nest("=", exp)

        if returnInit(exp[index[0]][index[1] - 1]) == returnInit(exp[index[0]][index[1] + 1]) {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "false")
        }
      } else {
        index := indexOf2Nest("!=", exp)

        if returnInit(exp[index[0]][index[1] - 1]) != returnInit(exp[index[0]][index[1] + 1]) {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "false")
        }
      }
    }

    for ;arrayContain2Nest(exp, "~~") || arrayContain2Nest(exp, "!~"); {
      if (indexOf2Nest("~~", exp)[0] < indexOf2Nest("!~", exp)[0] || indexOf2Nest("!~", exp)[0] == -1) && indexOf2Nest("!~", exp)[0] != -1 {
        index := indexOf2Nest("~~", exp)

        if isLess(add(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1], calc_params, line, functions), exp[index[0]][index[1] + 3]) || isLess(subtract(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1], calc_params, line, functions), exp[index[0]][index[1] + 3]) {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 3, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 3, "false")
        }
      } else {
        index := indexOf2Nest("!~", exp)

        if !(isLess(add(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1], calc_params, line, functions), exp[index[0]][index[1] + 4]) || isLess(subtract(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1], calc_params, line, functions), exp[index[0]][index[1] + 4])) {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 3, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 3, "false")
        }
      }
    }

    for ;arrayContain2Nest(exp, "~~~") || arrayContain2Nest(exp, "!~~"); {
      if (indexOf2Nest("~~~", exp)[0] < indexOf2Nest("!~~", exp)[0] || indexOf2Nest("!~~", exp)[0] == -1) && indexOf2Nest("~~~", exp)[0] != -1 {
        index := indexOf2Nest("~~~", exp)

        if returnInit(add(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 3], calc_params, line, functions)) == returnInit(exp[index[0]][index[1] + 3]) && returnInit(subtract(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1], calc_params, line, functions)) == returnInit(exp[index[0]][index[1] + 3]) {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 3, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 3, "false")
        }
      } else {
        index := indexOf2Nest("!~~", exp)

        if !(returnInit(add(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1], calc_params, line, functions)) == returnInit(exp[index[0]][index[1] + 3]) && returnInit(subtract(exp[index[0]][index[1] - 1], exp[index[0]][index[1] + 1], calc_params, line, functions)) == returnInit(exp[index[0]][index[1] + 3])) {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 3, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 3, "false")
        }
      }
    }

    for ;arrayContain2Nest(exp, "!"); {
      index := indexOf2Nest("!", exp)

      var operators_ = []string{"^", "*", "/", "%", "+", "-", "&", "|", ">", "<", ">=", "<=", "="}

      exp_ := []string{}

      for i := index[1] + 1; i < len(exp[index[0]]); i++ {
        if arrayContain(operators_, exp[index[0]][i]) {
          break
        }

        exp_ = append(exp_, exp[index[0]][i])
      }

      val := returnInit(mathParse(&[][]string{ exp_}, functions, line, calc_params, vars, dir)[0][0])

      if val == "false" || val == "undefined" || val == "null" {
        exp[index[0]] = replaceFrom(exp[index[0]], index[1], index[1] + len(exp_), "true")
      } else {
        exp[index[0]] = replaceFrom(exp[index[0]], index[1], index[1] + len(exp_), "false")
      }
    }

    for ;arrayContain2Nest(exp, "|") || arrayContain2Nest(exp, "&"); {
      if (indexOf2Nest("|", exp)[0] < indexOf2Nest("&", exp)[0] || indexOf2Nest("&", exp)[0] == -1) && indexOf2Nest("|", exp)[0] != -1 {
        index := indexOf2Nest("|", exp)

        val1 := returnInit(exp[index[0]][index[1] - 1])
        val2 := returnInit(exp[index[0]][index[1] + 1])

        if (val2 != "false" && val2 != "undefined" && val2 != "null") || (val1 != "false" && val1 != "undefined" && val1 != "null") {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "false")
        }
      } else {
        index := indexOf2Nest("&", exp)

        val1 := returnInit(exp[index[0]][index[1] - 1])
        val2 := returnInit(exp[index[0]][index[1] + 1])

        if (val2 != "false" && val2 != "undefined" && val2 != "null") && (val1 != "false" && val1 != "undefined" && val1 != "null") {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "true")
        } else {
          exp[index[0]] = replaceFrom(exp[index[0]], index[1] - 1, index[1] + 1, "false")
        }
      }
    }

    _exp := exp

    var copyer = [][]string{[]string{}}

    for i := 0; i < len(orig); i++ {
      for o := 0; o < len(orig[i]); o++ {
        copyer[len(copyer) - 1] = append(copyer[len(copyer) - 1], string(orig[i][o]))
      }
    }

    *gexp = [][]string(copyer)

    return _exp
  }
}
