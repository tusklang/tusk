package compiler

import . "lang/types"

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

        var statements = []string{ "local", "global", "log", "print", "cond", "while", "each", "include" } //list of statements

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
                  compilerErr("Conditionals require an array to loop through", dir, v.Line)
                }
            }

            if val == "include" { //if it is include, it is different than the other statements
              if right[0].Type != "string" {
                compilerErr("Expected a string after \"include\"", dir, v.Line)
              }

              included := includer(right[0].Value.(OmmString).ToGoType(), v.Line, dir)

              for _, acts := range included {
                actions = append(actions, acts...)
              }

            } else {
              actions = append(actions, Action{
                Type: val,
                Name: name,
                ExpAct: right,
              })
            }

            hasStatement = true
          }
        }

        if !hasStatement {
          compilerErr("\"" + (*v.Left).Item.Token.Name + "\" is not a statement", dir, v.Line)
        }
      case ":":

        if len(left) == 0 || left[0].Type != "variable" {
          compilerErr("Must have a variable before an assigner operator", dir, v.Line)
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
          compilerErr("Must have a variable before an increment or decrement", dir, v.Line)
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
          compilerErr("Must have a variable before an assignment operator", dir, v.Line)
        }
        if len(right) == 0 {
          compilerErr("Could not find a value after " + v.Type, dir, v.Line)
        }

        actions = append(actions, Action{
          Type: v.Type,
          Name: left[0].Name,
          Second: right,
        })

      case "cb-ob":
        //this is the operator to connect a closing brace to an opening brace
        //  like this:
        //  while (true) {}
        //between ) and { there must be an operator because everything in omm is an operation

        actions = append(actions, Action{
          Type: "cb-ob",
          First: left,
          ExpAct: right,
        })

      case "none":
        vActs := valueActions(v.Item, dir)

        actions = append(actions, vActs...)
    }
  }

  return actions
}
