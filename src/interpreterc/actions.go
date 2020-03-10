package main

import "strings"
import "regexp"
import "strconv"

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
  Indexes     [][]string
  ID              int
}

func actionizer(lex []string) []Action {
  var actions = []Action{}
  var len_lex = len(lex)

  actionReader:
  for i := 0; i < len_lex; i++ {

    switch (lex[i]) {
      case "newlineN":
        actions = append(actions, Action{ "newline", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 0 })
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

        actions = append(actions, Action{ "local", lex[i + 2], []string{}, exp, []string{}, []Action{}, []Condition{}, [][]string{}, 1 })
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

        actions = append(actions, Action{ "dynamic", lex[i + 2], []string{}, exp, []string{}, []Action{}, []Condition{}, [][]string{}, 2 })
        i+=(4 + len(exp_))
      case "alt":

        var alter = Action{ "alt", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 3 }

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

        actions = append(actions, Action{ "global", lex[i + 2], []string{}, exp, []string{}, []Action{}, []Condition{}, [][]string{}, 4 })
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

        actions = append(actions, Action{ "log", "", exp_, exp, []string{}, []Action{}, []Condition{}, [][]string{}, 5 })
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

        actions = append(actions, Action{ "print", "", exp_, exp, []string{}, []Action{}, []Condition{}, [][]string{}, 6 })
        i+=(2 + len(exp_))
      case "(":
        var exp = []string{}

        pCnt := 0

        for o := i; o < len_lex; o++ {

          if lex[o] == "(" {
            pCnt++
          }
          if lex[o] == ")" {
            pCnt--
          }

          if pCnt == 0 && lex[o] == "newlineS" {
            break
          }

          exp = append(exp, lex[o])
        }
        i+=len(exp)

        if i >= len_lex {
          actions = append(actions, Action{ "expression", "", exp, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 7 })
          break actionReader
        }

        if lex[i] == "." {
          indexes := [][]string{[]string{}}

          cbCnt := 0
          glCnt := 0
          bCnt := 0
          pCnt := 0

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

          i+=3

          actions = append(actions, Action{ "expressionIndex", "", exp, []Action{}, []string{}, []Action{}, []Condition{}, indexes, 8 })
        } else {
          actions = append(actions, Action{ "expression", "", exp, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 7 })
        }
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

        actions = append(actions, Action{ "group", "", []string{}, exp, []string{}, []Action{}, []Condition{}, [][]string{}, 9 })
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

          i+=len(logic_) + 1

          logic := actionizer(logic_)

          actions = append(actions, Action{ "process", procName, []string{}, logic, params, []Action{}, []Condition{}, [][]string{}, 10 })
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

          i+=len(logic_) + 1

          logic := actionizer(logic_)

          actions = append(actions, Action{ "process", "", []string{}, logic, params, []Action{}, []Condition{}, [][]string{}, 10 })
        }
        i--
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

        actions = append(actions, Action{ "#", name, []string{}, []Action{}, []string{}, params_, []Condition{}, [][]string{}, 11 })
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

        actions = append(actions, Action{ "return", "", []string{}, returner, []string{}, []Action{}, []Condition{}, [][]string{}, 12 })
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

        actions = append(actions, Action{ "conditional", "", []string{}, []Action{}, []string{}, []Action{}, conditions, [][]string{}, 13 })
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
        actions = append(actions, Action{ "import", "", []string{}, actionizedFile, []string{}, []Action{}, []Condition{}, [][]string{}, 14 })
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
        actions = append(actions, Action{ "read", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, [][]string{}, 15 })
        i+=(2 + len(phrase))
      case "break":
        actions = append(actions, Action{ "break", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 16 })
      case "skip":
        actions = append(actions, Action{ "skip", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 17 })
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
        actions = append(actions, Action{ "eval", "", []string{}, actionized, []string{}, []Action{}, []Condition{}, [][]string{}, 18 })
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
        actions = append(actions, Action{ "typeof", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, [][]string{}, 19 })
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
        actions = append(actions, Action{ "err", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, [][]string{}, 20 })
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

        actions = append(actions, Action{ "loop", "", []string{}, action, []string{}, []Action{}, []Condition{ { "loop", condition, action } }, [][]string{}, 21 })
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

        if i >= len_lex {
          actions = append(actions, Action{ "hash", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 22 })
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

          i+=3

          actions = append(actions, Action{ "hashIndex", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, indexes, 23 })
        } else {
          actions = append(actions, Action{ "hash", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 22 })
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

        if i >= len_lex {
          actions = append(actions, Action{ "array", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 24 })
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

          i+=3

          actions = append(actions, Action{ "arrayIndex", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, indexes, 25 })
        } else {
          actions = append(actions, Action{ "array", "", phrase, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 24 })
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
        actions = append(actions, Action{ "ascii", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, [][]string{}, 26 })
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
        actions = append(actions, Action{ "parse", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{}, [][]string{}, 27 })
        i+=(2 + len(phrase))
      default:

        if i < len_lex {

          if i + 1 < len_lex {

            if (lex[i + 1] == "++" || lex[i + 1] == "--") && strings.HasPrefix(lex[i], "$") {

              id_ := []byte(lex[i + 1])

              id := ""

              for o := 0; o < len(id_); o++ {
                _id := strconv.Itoa(int(id_[o]))
                id+=_id
              }

              intID, _ := strconv.Atoi(id)

              actions = append(actions, Action{ lex[i + 1], lex[i], []string{}, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, intID })
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

              actions = append(actions, Action{ lex[i + 1], lex[i], []string{}, by, []string{}, []Action{}, []Condition{}, [][]string{}, intID })
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

              actions = append(actions, Action{ "let", lex[i], exp_, exp, []string{}, []Action{}, []Condition{}, [][]string{}, 28 })
              i+=(len(exp))
            } else {

              var operators_ = []string{"^", "*", "/", "%", "+", "-", "&", "|", "!", ">", "<", ">=", "<=", "="}

              numMatch_, _ := regexp.MatchString("(\\d|\\.)+", lex[i])

              if !arrayContain(operators_, lex[i]) && !numMatch_ && !strings.HasPrefix(lex[i], "$") && strings.HasPrefix(lex[i], "'") && strings.HasPrefix(lex[i], "\"") && strings.HasPrefix(lex[i], "`") && lex[i] != "true" && lex[i] != "false" {
                break actionReader
              }

              var exp = []string{}

              pCnt := 0

              for o := i; o < len_lex; o++ {

                if !(o + 1 >= len_lex) {

                  if lex[o + 1] == "[" || lex[o + 1] == "]" {
                    break
                  }
                }

                if lex[o] == "(" {
                  pCnt++
                }
                if lex[o] == ")" {
                  pCnt--
                }

                numMatch, _ := regexp.MatchString("(\\d|\\.)+", lex[o])

                if !arrayContain(operators_, lex[o]) && !numMatch && !strings.HasPrefix(lex[o], "$") && strings.HasPrefix(lex[o], "'") && strings.HasPrefix(lex[o], "\"") && strings.HasPrefix(lex[o], "`") && lex[o] != "true" && lex[o] != "false" && pCnt == 0 {
                  exp = exp[:len(exp) - 1]
                  break
                }

                if pCnt == 0 && lex[o] == "newlineS" {
                  break
                }

                exp = append(exp, lex[o])
              }
              i+=len(exp)

              if i >= len_lex {
                actions = append(actions, Action{ "expression", "", exp, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 7 })
                break actionReader
              }

              if lex[i] == "." {

                indexes := [][]string{[]string{}}

                cbCnt := 0
                glCnt := 0
                bCnt := 0
                pCnt := 0

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

                i+=3

                actions = append(actions, Action{ "expressionIndex", "", exp, []Action{}, []string{}, []Action{}, []Condition{}, indexes, 8 })
              } else {
                actions = append(actions, Action{ "expression", "", exp, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 7 })
              }
            }
          } else {

            var operators_ = []string{"^", "*", "/", "%", "+", "-", "&", "|", "!", ">", "<", ">=", "<=", "="}

            numMatch_, _ := regexp.MatchString("(\\d|\\.)+", lex[i])

            if !arrayContain(operators_, lex[i]) && !numMatch_ && !strings.HasPrefix(lex[i], "$") && strings.HasPrefix(lex[i], "'") && strings.HasPrefix(lex[i], "\"") && strings.HasPrefix(lex[i], "`") && lex[i] != "true" && lex[i] != "false" {
              continue
            }

            var exp = []string{}

            pCnt := 0

            for o := i; o < len_lex; o++ {

              if !(o + 1 >= len_lex) {

                if lex[o + 1] == "[" || lex[o + 1] == "]" {
                  break
                }
              }

              if lex[o] == "(" {
                pCnt++
              }
              if lex[o] == ")" {
                pCnt--
              }

              numMatch, _ := regexp.MatchString("(\\d|\\.)+", lex[o])

              if !arrayContain(operators_, lex[o]) && !numMatch && !strings.HasPrefix(lex[o], "$") && strings.HasPrefix(lex[o], "'") && strings.HasPrefix(lex[o], "\"") && strings.HasPrefix(lex[o], "`") && lex[o] != "true" && lex[o] != "false" && pCnt == 0 {
                exp = exp[:len(exp) - 1]
                break
              }

              if pCnt == 0 && lex[o] == "newlineS" {
                break
              }

              exp = append(exp, lex[o])
            }
            i+=len(exp)

            if i >= len_lex {
              actions = append(actions, Action{ "expression", "", exp, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 7 })
              break actionReader
            }

            if lex[i] == "." {

              indexes := [][]string{[]string{}}

              cbCnt := 0
              glCnt := 0
              bCnt := 0
              pCnt := 0

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

              i+=3

              actions = append(actions, Action{ "expressionIndex", "", exp, []Action{}, []string{}, []Action{}, []Condition{}, indexes, 8 })
            } else {
              actions = append(actions, Action{ "expression", "", exp, []Action{}, []string{}, []Action{}, []Condition{}, [][]string{}, 7 })
            }
          }
        }
    }
  }

  return actions
}
