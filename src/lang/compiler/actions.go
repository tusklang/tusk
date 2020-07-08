package compiler

import "fmt"
import "os"
import "strconv"

import . "lang/types"

func actionizerErr(msg string, dir string, line uint64) {
  colorprint("Error while actionizing " + dir + " at line " + strconv.FormatUint(line, 10) + "\n", 12)
  fmt.Println(msg)
  os.Exit(1)
}

func actionizer(operations []Operation, dir string) []Action {

  var actions []Action

  for _, v := range operations {

    var left []Action
    var right []Action

    if v.Type != "none" {

      if v.Left != nil {
        left = actionizer([]Operation{ *v.Left }, dir)
      }
      if v.Right != nil {
        right = actionizer([]Operation{ *v.Right }, dir)
      }

    }

    switch v.Type {
      case "~":

        var statements = []string{ "local", "global", "log", "print", "cond", "while", "each" } //list of statements

        var hasStatement bool = false

        for _, val := range statements {
          if val == (*v.Left).Item.Token.Name {

            var name string

            if len(right) > 0 {
              name = right[0].Name
            }

            switch val {
              case "cond":
                //if the cond does not have an array
                //like this
                //  cond [
                //    true ? {_},
                //    false ? {_},
                //  ]
                //give an error

                if right[0].Type != "array" {
                  actionizerErr("Conditionals require an array to loop through", dir, v.Line)
                }
            }

            actions = append(actions, Action{
              Type: val,
              Name: name,
              ExpAct: right,
            })

            hasStatement = true
          }
        }

        if !hasStatement {
          actionizerErr("\"" + (*v.Left).Item.Token.Name + "\" is not a statement", dir, v.Line)
        }
      case ":":

        if len(left) == 0 || left[0].Type != "variable" {
          actionizerErr("Must have a variable before an assigner operator", dir, v.Line)
        }

        varname := left[0].Name
        value := right

        actions = append(actions, Action{
          Type: "let",
          Name: varname,
          ExpAct: value,
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
      case "~~": fallthrough
      case "~~~": fallthrough
      case "!": fallthrough
      case "&": fallthrough
      case "|": fallthrough
      case "::": fallthrough
      case "?": fallthrough
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
        })
      ////////////////////////////////////////////////////////

      case "++": fallthrough
      case "--":

        if len(left) == 0 || left[0].Type != "variable" {
          actionizerErr("Must have a variable before an increment or decrement", dir, v.Line)
        }

        actions = append(actions, Action{
          Type: v.Type,
          Name: left[0].Name,
        })

      case "+=": fallthrough
      case "-=": fallthrough
      case "*=": fallthrough
      case "/=": fallthrough
      case "%=": fallthrough
      case "^=":

        if len(left) == 0 || left[0].Type != "variable" {
          actionizerErr("Must have a variable before an assignment operator", dir, v.Line)
        }
        if len(right) == 0 {
          actionizerErr("Could not find a value after " + v.Type, dir, v.Line)
        }

        actions = append(actions, Action{
          Type: v.Type,
          Name: left[0].Name,
          Second: right,
        })

      case "none":
        vActs := valueActions(v.Item, dir)

        actions = append(actions, vActs...)
    }
  }

  return actions
}
