package main

import "fmt"
import "os"
import "strings"
import "strconv"
import "bufio"
import "regexp"
import "os/exec"

//json lib created by the jsoniter group on github
import "github.com/json-iterator/go"

type Variable struct {
  Name        string
  Type        string
  Value     []string
  ValueActs []Action
}

type Returner struct {
  Val       []string
  Variables   map[string]Variable
  Exp     [][]string
  Type        string
}

func mergeVars(map1, map2 map[string]Variable) map[string]Variable {
  combined := map[string]Variable{}

  for k, v := range map1 {
    combined[k] = v
  }

  for k, v := range map2 {
    combined[k] = v
  }

  return combined;
}

func parser(actions []Action, calc_params paramCalcOpts, dir string, line_ uint64, functions []Funcs, vars_ map[string]Variable, groupReturn bool) Returner {

  var json = jsoniter.ConfigCompatibleWithStandardLibrary

  var vars = vars_
  var line = line_

  scanner := bufio.NewReader(os.Stdin)

  var exp = [][]string{}

  for i := 0; i < len(actions); i++ {

    switch (actions[i].Type) {
      case "newlineN":
        exp = [][]string{}
        line++
      case "let":
        vars[actions[i].Name] = Variable{ actions[i].Name, "local", parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false).Exp[0], []Action{} }
      case "abstract":
        vars[actions[i].Name] = Variable{ actions[i].Name, "abstract", []string{}, actions[i].ExpAct }
      case "alt":

        o := 0

        for ;parser(actions[i].Condition[0].Condition, calc_params, dir, line, functions, vars, false).Exp[0][0] == "true"; {
          if o >= len(actions[i].Condition) {
            o = 0
          }

          parsed := parser(actions[i].Condition[o].Actions, calc_params, dir, line, functions, vars, false)

          for k, v := range parsed.Variables {
            vars[k] = v
          }

          o++
        }
      case "global":
        vars[actions[i].Name] = Variable{ actions[i].Name, "global", parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false).Exp[0], []Action{} }
      case "log":
        log(strings.Join(parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false).Exp[0], ""))
      case "print":
        log(strings.Join(parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false).Exp[0], ""))
      case "expression":
        expStr := [][]string{ actions[i].ExpStr }
        exp = append(exp, mathParse(&expStr, functions, line, calc_params, vars, dir)...)
        actions[i].ExpStr = expStr[0]
      case "group":
        grouped := parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false)

        if groupReturn {
          return grouped
        }

        if grouped.Type == "return" || grouped.Type == "break" || grouped.Type == "skip" {
          return Returner{ grouped.Val, grouped.Variables, grouped.Exp, grouped.Type }
        }

      case "process":
        if actions[i].Name == "" {
          fmt.Println("There Was An Error: Anonymous Processes Are Not Yet Supported\n\nproc(\n^^^^ <-- Error On Line " + strconv.Itoa(int(line)))
          os.Exit(1)
        } else {
          vars[actions[i].Name] = Variable{ actions[i].Name, "process", actions[i].Params, actions[i].ExpAct }
        }
      case "#":
        proc := vars[actions[i].Name]

        if proc.Name == "" {
          fmt.Println("There Was An Error: " + actions[i].Name + " is not a function\n\n#~" + actions[i].Name[1:] + "(\n^^^ <-- Error On Line " + strconv.Itoa(int(line)))
          os.Exit(1)
        }

        procParams := proc.Value
        nParams := make(map[string]Variable)

        for o := 0; o < len(procParams); o++ {

          if len(actions[i].Args) <= o {
            break
          }

          nParams[procParams[o]] = Variable{ procParams[o], "local", parser([]Action{ actions[i].Args[o] }, calc_params, dir, line, functions, vars, true).Exp[0], []Action{} }
        }

        params := make(map[string]Variable)

        for k, v := range vars {
          params[k] = v
        }

        for k, v := range nParams {
          params[k] = v
        }

        parsed := parser(proc.ValueActs, calc_params, dir, line, functions, params, true)

        if len(parsed.Val) <= 0 {
          exp = append(exp, []string{ "undefined" })
        } else {
          val := [][]string{ parsed.Val }
          exp = append(exp, mathParse(&val, functions, line, calc_params, vars, dir)...)
          parsed.Val = val[0]
        }
      case "return":
        return Returner{ parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false).Exp[0], vars, mathParse(&exp, functions, line, calc_params, vars, dir), "return" }
      case "conditional":

        for o := 0; o < len(actions[i].Condition); o++ {

          opars := parser(actions[i].Condition[o].Condition, calc_params, dir, line, functions, vars, false).Exp

          val := mathParse(&opars, functions, line, calc_params, vars, dir)[0][0]

          if val != "false" && val != "undefined" && val != "null" {
            parsed := parser(actions[i].Condition[o].Actions, calc_params, dir, line, functions, vars, false)

            vars = mergeVars(vars, parsed.Variables)

            if parsed.Type == "return" {
              return Returner{ parsed.Val, parsed.Variables, parsed.Exp, "return" }
            }

            if parsed.Type == "break" {
              return Returner{ parsed.Val, parsed.Variables, parsed.Exp, "break" }
            }

            if parsed.Type == "skip" {
              return Returner{ parsed.Val, parsed.Variables, parsed.Exp, "skip" }
            }

            break
          }
        }
      case "import":

        fileName := parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false).Exp[0][0]

        file := read(dir + fileName, "Cannot Find File: " + fileName, true)
        fileNQ, _ := NQReplace(file)

        lexCmd := exec.Command("./lexer/main-win.exe")

        lexCmd.Stdin = strings.NewReader(fileNQ + "\n")

        _lex, _ := lexCmd.CombinedOutput()
        lex_ := string(_lex)

        var lex []string

        json.Unmarshal([]byte(lex_), &lex)

        actions := actionizer(lex)
        parsed := parser(actions, calc_params, dir, line, functions, vars, false)

        for k, v := range parsed.Variables {
          vars[k] = v
        }
      case "hash":
        exp = append(exp, actions[i].ExpStr)
      case "array":
        exp = append(exp, actions[i].ExpStr)
      case "hashIndex":

        val, _ := hashIndex(actions[i].ExpStr, actions[i].Indexes, functions, line, calc_params, vars, dir)

        exp = append(exp, mathParse(&[][]string{ val }, functions, line, calc_params, vars, dir)[0])
      case "arrayIndex":

        val, _ := arrayIndex(actions[i].ExpStr, actions[i].Indexes, functions, line, calc_params, vars, dir)

        exp = append(exp, mathParse(&[][]string{ val }, functions, line, calc_params, vars, dir)[0])
      case "expressionIndex":

        var val []string

        expStr := mathParse(&[][]string{ actions[i].ExpStr }, functions, line, calc_params, vars, dir)[0]

        if expStr[0] == "[:" {
          val, _ = hashIndex(expStr, actions[i].Indexes, functions, line, calc_params, vars, dir)
        } else if expStr[0] == "[" {
          val, _ = arrayIndex(expStr, actions[i].Indexes, functions, line, calc_params, vars, dir)
        } else {
          val, _ = stringIndex(expStr, actions[i].Indexes, functions, line, calc_params, vars, dir)
        }

        exp = append(exp, mathParse(&[][]string{ val }, functions, line, calc_params, vars, dir)[0])
      case "read":
        log(strings.Join(parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false).Exp[0], ""))
        text, _ := scanner.ReadString('\n')

        exp = append(exp, []string{ text })
      case "break":
        return Returner{ []string{"break"}, vars, [][]string{}, "break" }
      case "skip":
        return Returner{ []string{"skip"}, vars, [][]string{}, "skip" }
      case "eval":

        calculated := mathParse(&[][]string{actions[i].ExpStr}, functions, line, calc_params, vars, dir)[0][0]

        fileNQ, _ := NQReplace(calculated)

        lexCmd := exec.Command("./lexer/main-win.exe")

        lexCmd.Stdin = strings.NewReader(fileNQ + "\n")

        _lex, _ := lexCmd.CombinedOutput()
        lex_ := string(_lex)

        var lex []string

        json.Unmarshal([]byte(lex_), &lex)

        actions := actionizer(lex)
        parser(actions, calc_params, dir, line, functions, vars, false)
      case "typeof":

        parsed := parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false).Exp[0][0]

        numCheck, _ := regexp.MatchString("\\d+", parsed)

        if strings.HasPrefix(parsed, "'") || strings.HasPrefix(parsed, "\"") || strings.HasPrefix(parsed, "`") {
          exp = append(exp, []string{ "string" })
        } else if strings.HasPrefix(parsed, "[:") {
          exp = append(exp, []string{ "hash" })
        } else if strings.HasPrefix(parsed, "[") {
          exp = append(exp, []string{ "array" })
        } else if parsed == "true" || parsed == "false" {
          exp = append(exp, []string{ "boolean" })
        } else if parsed == "undefined" || parsed == "null" {
          exp = append(exp, []string{ "falsey" })
        } else if numCheck {
          exp = append(exp, []string{ "number" })
        }
      case "err":
        parsed := parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false).Exp[0][0]
        fmt.Println("There Was An Error: " + parsed + " On Line " + strconv.Itoa(int(line)))
      case "loop":

        cond := [][]string{ actions[i].Condition[0].Condition[0].ExpStr }

        curBool := mathParse(&cond, functions, line, calc_params, vars, dir)[0][0]

        for ;curBool != "false" && curBool != "undefined" && curBool != "null"; {
          var parsed = parser(actions[i].Condition[0].Actions, calc_params, dir, line, functions, vars, false)

          vars = mergeVars(vars, parsed.Variables)

          curBool = mathParse(&cond, functions, line, calc_params, vars, dir)[0][0]

          if parsed.Type == "return" {
            return Returner{ parsed.Val, parsed.Variables, parsed.Exp, "return" }
          }

          if parsed.Type == "break" {
            break
          }

          if parsed.Type == "skip" {
            continue
          }
        }
      case "ascii":
        parsed := parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false)
        calculated := mathParse(&[][]string{ parsed.Exp[0] }, functions, line, calc_params, vars, dir)[0][0]

        if !(strings.HasPrefix(calculated, "'") || strings.HasPrefix(calculated, "\"") || strings.HasPrefix(calculated, "`")) {
          fmt.Println("There Was An Error: You Cannot Get An ASCII Value Of A Non-String\n\n" + calculated + "\n^ <- Error On Line " + strconv.Itoa(int(line)))
          os.Exit(1)
        }

        calculated_ := []rune(calculated[1:][:len(calculated) - 2])[0]
        exp = append(exp, []string{ strconv.Itoa(int(calculated_)) })
      case "parse":
        parsed := parser(actions[i].ExpAct, calc_params, dir, line, functions, vars, false)
        calculated := mathParse(&[][]string{ parsed.Exp[0] }, functions, line, calc_params, vars, dir)[0][0]

        numCheck, _ := regexp.MatchString("\\d+", calculated)

        if !(strings.HasPrefix(calculated, "'") || strings.HasPrefix(calculated, "\"") || strings.HasPrefix(calculated, "`")) && !numCheck {
          fmt.Println("There Was An Error: You Cannot Parse A Non-String Or Non-Number Value")
          os.Exit(1)
        }

        if strings.HasPrefix(calculated, "'") || strings.HasPrefix(calculated, "\"") || strings.HasPrefix(calculated, "`") {
          exp = append(exp, []string{ calculated[1:len(calculated) - 1] })
        } else {
          exp = append(exp, []string{ calculated })
        }
    }
  }

  return Returner{ []string{}, vars, mathParse(&exp, functions, line, calc_params, vars, dir), "" }
}
