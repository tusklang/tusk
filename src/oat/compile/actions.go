package main

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
}

func actionizer(lex []string) []Action {
  var actions = []Action{}
  var len_lex = len(lex)

  for i := 0; i < len_lex; i++ {
    switch (lex[i]) {
      case "newlineN":
        actions = append(actions, Action{ "newline", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{} })
        continue;
      case "let":
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

        actions = append(actions, Action{ "let", lex[i + 2], exp_, exp, []string{}, []Action{}, []Condition{} })
        i+=(4 + len(exp_))
      case "abstract":
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

        actions = append(actions, Action{ "abstract", lex[i + 2], exp_, exp, []string{}, []Action{}, []Condition{} })
        i+=(4 + len(exp_))
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

        actions = append(actions, Action{ "let", lex[i + 2], exp_, exp, []string{}, []Action{}, []Condition{} })
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

        actions = append(actions, Action{ "log", "", exp_, exp, []string{}, []Action{}, []Condition{} })
        i+=(2 + len(exp_))
      case "[:":
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

          if glCnt == 0 {
            exp_ = append(exp_, lex[o])
            break
          }

          exp_ = append(exp_, lex[o])
        }

        actions = append(actions, Action{ "glossary", "", exp_, []Action{}, []string{}, []Action{}, []Condition{} })
        i+=len(exp_)
      case "[":
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

          if bCnt == 0 {
            exp_ = append(exp_, lex[o])
            break;
          }

          exp_ = append(exp_, lex[o])
        }

        exp := actionizer(exp_)
        actions = append(actions, Action{ "array", "", exp_, exp, []string{}, []Action{}, []Condition{} })
        i+=len(exp_)
      case "(":
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

          if pCnt == 0 {
            exp_ = append(exp_, lex[o])
            break
          }

          exp_ = append(exp_, lex[o])
        }

        actions = append(actions, Action{ "expression", "", exp_, []Action{}, []string{}, []Action{}, []Condition{} })
        i+=len(exp_)
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

          if cbCnt == 0 {
            break
          }

          if lex[o] == "newlineS" {
            break
          }

          exp_ = append(exp_, lex[o])
        }

        exp_ = exp_[1:]
        exp_ = exp_[:1]

        exp := actionizer(exp_)

        actions = append(actions, Action{ "group", "", exp_, exp, []string{}, []Action{}, []Condition{} })
        i+=(len(exp_))
      case "process":
        if lex[i + 1] == "~" {
          var procName = lex[i + 2]
          params := []string{}

          for o := i + 4; o < len_lex; o+=2 {
            if lex[o] == ")" {
              break
            }

            params = append(params, lex[o])
          }
          i+=(5 + len(params))

          var logic_ = []string{}

          cbCnt := 0

          for o := i; o < len_lex; o++ {
            if lex[o] == "{" {
              cbCnt++
            }

            if lex[o] == "}" {
              cbCnt--
            }

            if cbCnt == 0 && lex[o] == "newlineS" {
              break
            }

            logic_ = append(logic_, lex[o])
          }

          logic := actionizer(logic_)

          actions = append(actions, Action{ "process", procName, []string{}, logic, params, []Action{}, []Condition{} })
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

            if cbCnt == 0 && lex[o] == "newlineS" {
              break
            }

            logic_ = append(logic_, lex[o])
          }

          logic := actionizer(logic_)

          actions = append(actions, Action{ "process", "", []string{}, logic, params, []Action{}, []Condition{} })
        }
      case "#":

        cbCnt := 0
        glCnt := 0
        bCnt := 0
        pCnt := 1

        var name = lex[i + 1]

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

        actions = append(actions, Action{ "#", name, []string{}, []Action{}, []string{}, params_, []Condition{} })
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

        actions = append(actions, Action{ "return", "", []string{}, returner, []string{}, []Action{}, []Condition{}})
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

        actions = append(actions, Action{ "conditional", "", []string{}, []Action{}, []string{}, []Action{}, conditions })
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
        actions = append(actions, Action{ "import", "", []string{}, actionizedFile, []string{}, []Action{}, []Condition{} })
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
        actions = append(actions, Action{ "read", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{} })
        i+=(2 + len(phrase))
      case "break":
        actions = append(actions, Action{ "break", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{} })
      case "skip":
        actions = append(actions, Action{ "skip", "", []string{}, []Action{}, []string{}, []Action{}, []Condition{} })
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

        actions = append(actions, Action{ "eval", "", phrase, []Action{}, []string{}, []Action{}, []Condition{} })
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
        actions = append(actions, Action{ "typeof", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{} })
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
        actions = append(actions, Action{ "err", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{} })
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

        actions = append(actions, Action{ "loop", "", []string{}, action, []string{}, []Action{}, []Condition{ { "loop", condition, action } } })
        i+=(1 + len(condition_) + len(action_))
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
        actions = append(actions, Action{ "ascii", "", []string{}, actionizedPhrase, []string{}, []Action{}, []Condition{} })
        i+=(2 + len(phrase))
      default:

        if i < len_lex {

          if i + 1 < len_lex {
            if lex[i + 1] == ":" {
              exp_ := []string{}

              for o := i; o < len_lex; o++ {

                if lex[o] == "newlineS" {
                  break;
                }

                exp_ = append(exp_, lex[o]);
              }

              exp := actionizer(exp_)

              actions = append(actions, Action{ "let", "", exp_, exp, []string{}, []Action{}, []Condition{} })
              i+=(len(exp))
            }
          } else {
            var exp []string

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

              if cbCnt == 0 && glCnt == 0 && bCnt == 0 && pCnt == 0 && lex[o] == "newlineS" {
                break
              }

              exp = append(exp, lex[o])
            }

            actions = append(actions, Action{ "expression", "", exp, []Action{}, []string{}, []Action{}, []Condition{} })
            i+=len(exp)
          }
        }
    }
  }

  return actions
}
