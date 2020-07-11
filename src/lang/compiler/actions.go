package compiler

import . "lang/types"

func actionizer(operations []Operation, dir string) []Action {

  var actions []Action

  for _, v := range operations {

    var left []Action
    var right []Action

    if v.Left != nil {
      left = actionizer([]Operation{ *v.Left }, dir)
    }
    if v.Right != nil {
      right = actionizer([]Operation{ *v.Right }, dir)
    }

    switch v.Type {
      case "~":

        var statements = []string{ "local", "global", "log", "print", "if", "elif", "else", "while", "each", "include", "function", "return" } //list of statements

        var hasStatement bool = false

        for _, val := range statements {
          if val == (*v.Left).Item.Token.Name {

            var name string

            if len(right) > 0 {
              name = right[0].Name

              if len(right[0].Indexes) != 0 {
                compilerErr("Cannot use the index operator (::) for a statement", dir, v.Line)
              }
            }

            switch val {

              case "include":
                if right[0].Type != "string" {
                  compilerErr("Expected a string after \"include\"", dir, v.Line)
                }

                includeFiles := includer(right[0].Value.(OmmString).ToGoType(), v.Line, dir)

                for _, acts := range includeFiles {
                  actions = append(actions, acts...)
                }

              case "function":
                if right[0].Type != "=>" {
                  compilerErr("Expected a => operator to connect the function parameter list and the function body", dir, right[0].Line)
                }

                for _, p := range right[0].First[0].ExpAct {
                  if p.Type == "global" || p.Type == "local" {
                    compilerErr("Cannot set access for parameter defaults", dir, right[0].Line)
                  }
                  if p.Type != "let" && p.Type != "variable" {
                    compilerErr("Function parameter lists can only have let statements and variables", dir, right[0].Line)
                  }
                }

                actions = append(actions, Action{
                  Type: "function",
                  Value: OmmFunc{
                    Params: right[0].First[0].ExpAct, //getting the ExpAct because it wont matter if the dev uses a { or a ( because everything will be in the function's scope
                    Body: right[0].Second,
                  },
                  File: dir,
                  Line: v.Line,
                })

              case "if":

                if right[0].Type != "=>" {
                  compilerErr("Expected a => operator to connect the if condition and the body", dir, right[0].Line)
                }

                actions = append(actions, Action{
                  Type: "condition",
                  ExpAct: []Action{ Action{
                    Type: "if",
                    First: right[0].First,
                    ExpAct: right[0].Second,
                  } },
                  File: dir,
                  Line: v.Line,
                })

              case "elif":

                if right[0].Type != "=>" {
                  compilerErr("Expected a => operator to connect the elif condition and the body", dir, right[0].Line)
                }

                if len(actions) == 0 || actions[len(actions) - 1].Type != "condition" {
                  compilerErr("Unexpected elif statement", dir, right[0].Line)
                }

                //append to the previous conditional statement
                actions[len(actions) - 1].ExpAct = append(actions[len(actions) - 1].ExpAct, Action{
                  Type: "if",
                  First: right[0].First,
                  ExpAct: right[0].Second,
                })

              case "else":

                if len(actions) == 0 || actions[len(actions) - 1].Type != "condition" {
                  compilerErr("Unexpected else statement", dir, right[0].Line)
                }

                //append to the previous conditional statement
                actions[len(actions) - 1].ExpAct = append(actions[len(actions) - 1].ExpAct, Action{
                  Type: "else",
                  ExpAct: right,
                })

              case "while":
                if right[0].Type != "=>" {
                  compilerErr("Expected a => operator to connect the while condition and the body", dir, right[0].Line)
                }

                actions = append(actions, Action{
                  Type: val,
                  First: right[0].First,
                  ExpAct: right[0].Second,
                  File: dir,
                  Line: v.Line,
                })

              case "each":
                if right[0].Type != "=>" {
                  compilerErr("Expected a => operator to connect the each iterator and the body", dir, right[0].Line)
                }

                if len(right[0].First[0].ExpAct) != 3 {
                  compilerErr("Each loops must look like this: each(iterator, key, value)", dir, right[0].Line)
                }

                for _, n := range right[0].First[0].ExpAct[1:] {
                  if n.Type != "variable" {
                    compilerErr("Key or value was not given as a variable", dir, right[0].Line)
                  }
                }

                actions = append(actions, Action{
                  Type: val,
                  First: right[0].First[0].ExpAct, //because it doesnt matter if they use a { or (
                  ExpAct: right[0].Second,
                  File: dir,
                  Line: v.Line,
                })

              default:

                actions = append(actions, Action{
                  Type: val,
                  Name: name,
                  ExpAct: right,
                  File: dir,
                  Line: v.Line,
                })

            }

            hasStatement = true
          }
        }

        if !hasStatement {
          compilerErr("\"" + (*v.Left).Item.Token.Name + "\" is not a statement", dir, v.Line)
        }
      case ":":

        if len(left) == 0 || (left[0].Type != "variable" && left[0].Type != "::")  {
          compilerErr("Must have a variable before an assigner operator", dir, v.Line)
        }

        var varname string
        var indexes [][]Action

        if left[0].Type == "variable" {
          varname = left[0].Name
        } else {

          var currentIndex = left[0]

          for ;currentIndex.Type != "variable"; {
            indexes = append([][]Action{ currentIndex.Second }, indexes...)
            currentIndex = currentIndex.First[0]
          }

          varname = currentIndex.Name
        }

        value := right

        actions = append(actions, Action{
          Type: "let",
          Name: varname,
          ExpAct: value,
          Indexes: indexes,
          File: dir,
          Line: v.Line,
        })

      //all of these operations have the same way of appending
      case "+": fallthrough
      case "-": fallthrough
      case "*": fallthrough
      case "/": fallthrough
      case "%": fallthrough
      case "^": fallthrough
      case "=": fallthrough
      case "!=": fallthrough
      case ">": fallthrough
      case "<": fallthrough
      case ">=": fallthrough
      case "<=": fallthrough
      case "!": fallthrough
      case "&": fallthrough
      case "|": fallthrough
      case "::": fallthrough
      case "->": fallthrough
      case "=>": fallthrough
      case "sync": fallthrough
      case "async":

        var degree []Action

        if v.Degree != nil {
          degree = actionizer([]Operation{ *v.Degree }, dir)
        }

        actions = append(actions, Action{
          Type: v.Type,
          First: left,
          Second: right,
          Degree: degree,
          File: dir,
          Line: v.Line,
        })
      ////////////////////////////////////////////////////////

      case "++": fallthrough
      case "--":

        if len(left) == 0 || (left[0].Type != "variable" && left[0].Type != "::") {
          compilerErr("Must have a variable before an increment or decrement", dir, v.Line)
        }

        var varname string
        var indexes [][]Action

        if left[0].Type == "variable" {
          varname = left[0].Name
        } else {

          var currentIndex = left[0]

          for ;currentIndex.Type != "variable"; {
            indexes = append([][]Action{ currentIndex.Second }, indexes...)
            currentIndex = currentIndex.First[0]
          }

          varname = currentIndex.Name
        }

        actions = append(actions, Action{
          Type: v.Type,
          Name: varname,
          File: dir,
          Line: v.Line,
        })

      case "+=": fallthrough
      case "-=": fallthrough
      case "*=": fallthrough
      case "/=": fallthrough
      case "%=": fallthrough
      case "^=":

        if len(left) == 0 || (left[0].Type != "variable" && left[0].Type != "::") {
          compilerErr("Must have a variable before an assignment operator", dir, v.Line)
        }
        if len(right) == 0 {
          compilerErr("Could not find a value after " + v.Type, dir, v.Line)
        }

        var varname string
        var indexes [][]Action

        if left[0].Type == "variable" {
          varname = left[0].Name
        } else {

          var currentIndex = left[0]

          for ;currentIndex.Type != "variable"; {
            indexes = append([][]Action{ currentIndex.Second }, indexes...)
            currentIndex = currentIndex.First[0]
          }

          varname = currentIndex.Name
        }

        actions = append(actions, Action{
          Type: v.Type,
          Name: varname,
          Second: right,
          File: dir,
          Line: v.Line,
        })

      case "break": fallthrough
      case "continue":

        actions = append(actions, Action{
          Type: v.Type,
          File: dir,
          Line: v.Line,
        })

      case "none":
        vActs := valueActions(v.Item, dir)

        actions = append(actions, vActs)
    }
  }

  return actions
}
