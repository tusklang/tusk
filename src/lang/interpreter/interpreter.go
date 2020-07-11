package interpreter

import "fmt"
import "os"
import . "lang/types"

func ommPanic(err string, line uint64, file string) {
  fmt.Println("Panic on line", line, "file", file)
  fmt.Println(err)
  os.Exit(1)
}

func interpreter(actions []Action, cli_params CliParams, passedVars map[string]Variable, this_vals []Action) Returner {

  var expReturn = false //if it is inside an expression

  if len(actions) == 1 {
    expReturn = true
  }

  var vars = make(map[string]Variable)

  for k, v := range passedVars { //copy passedVars into vars (so they won't mutate)
    vars[k] = v
  }

  for _, v := range actions {
    switch v.Type {

      case "local": fallthrough
      case "global": fallthrough
      case "let":

        interpreted := interpreter(v.ExpAct, cli_params, vars, this_vals)

        for k, variable := range interpreted.Variables {
          vars[k] = variable
        }

        vars[v.Name] = Variable{
          Type: v.Type,
          Value: interpreted.Exp,
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: vars[v.Name].Value,
            Variables: vars,
          }
        }

      case "log":
        interpreted := interpreter(v.ExpAct, cli_params, vars, this_vals)
        for k, variable := range interpreted.Variables {
          vars[k] = variable
        }
        fmt.Println(interpreted.Exp.Format())
      case "print":
        interpreted := interpreter(v.ExpAct, cli_params, vars, this_vals)
        for k, variable := range interpreted.Variables {
          vars[k] = variable
        }
        fmt.Print(interpreted.Exp.Format())

      //all of the types
      case "string": fallthrough
      case "rune": fallthrough
      case "number": fallthrough
      case "bool": fallthrough
      case "function": fallthrough
      case "array": fallthrough
      case "hash": fallthrough
      case "undef": fallthrough
      case "thread":

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: v.Value,
            Variables: vars,
          }
        }
      //////////////////

      case "variable":

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: vars[v.Name].Value,
            Variables: vars,
          }
        }

      case "{":

        groupRet := interpreter(v.ExpAct, cli_params, vars, this_vals)

        for k, variable := range groupRet.Variables {
          if _, exists := vars[k]; exists || variable.Type == "global" {
            vars[k] = variable
          }
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: groupRet.Exp,
            Variables: groupRet.Variables,
          }
        }

      case "(":

        groupRet := interpreter(v.ExpAct, cli_params, vars, this_vals)

        for k, variable := range groupRet.Variables {
          vars[k] = variable
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: groupRet.Exp,
            Variables: vars,
          }
        }

        //operations
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
        case "->": fallthrough
        case "=>": fallthrough //this is probably not necessary, but i just left it here
        case "sync": fallthrough
        case "async":

          firstInterpreted := interpreter(v.First, cli_params, vars, this_vals)
          secondInterpreted := interpreter(v.Second, cli_params, vars, this_vals)

          for k, variable := range firstInterpreted.Variables {
            vars[k] = variable
          }
          for k, variable := range secondInterpreted.Variables {
            vars[k] = variable
          }

          fmt.Println(v.Type)

          operationFunc, exists := operations[firstInterpreted.Exp.Type() + " " + v.Type + " " + secondInterpreted.Exp.Type()]

          if !exists { //if there is no operation for that type, panic
            ommPanic("Could not find " + v.Type + " operation for types " + firstInterpreted.Exp.Type() + " and " + secondInterpreted.Exp.Type(), v.Line, v.File)
          }

          computed := operationFunc(firstInterpreted.Exp, secondInterpreted.Exp, cli_params, v.Line, v.File)

          if expReturn {
            return Returner{
              Type: "expression",
              Exp: computed,
              Variables: vars,
            }
          }

        ////////////

        case "break": fallthrough
        case "continue":

          return Returner{
            Type: v.Type,
            Variables: vars,
          }

        case "return":

          return Returner{
            Type: "return",
            Exp: interpreter(v.ExpAct, cli_params, vars, this_vals).Exp,
            Variables: vars,
          }

        case "while":

          cond := interpreter(v.First, cli_params, vars, this_vals)

          for k, variable := range cond.Variables {
            vars[k] = variable
          }

          for ;isTruthy(cond.Exp); cond = interpreter(v.First, cli_params, vars, this_vals) {
            for k, variable := range cond.Variables {
              vars[k] = variable
            }

            interpreted := interpreter(v.ExpAct, cli_params, vars, this_vals)

            for k, variable := range interpreted.Variables {
              vars[k] = variable
            }

            if interpreted.Type == "return" {
              return Returner{
                Type: interpreted.Type,
                Exp: interpreted.Exp,
                Variables: vars,
              }
            }

            if interpreted.Type == "break" {
              break
            }
            if interpreted.Type == "continue" {
              continue;
            }
          }

        case "each":

          it := interpreter([]Action{ v.First[0] }, cli_params, vars, this_vals).Exp
          keyName := v.First[1].Name //get name of key
          valName := v.First[2].Name //get name of val

          switch it.(type) {
            case OmmHash:

              for key, val := range it.(OmmHash).Hash {

                ommtypeKey := OmmString{}
                ommtypeKey.FromGoType(key)

                var sendVars = make(map[string]Variable)

                for k, variable := range vars {
                  sendVars[k] = variable
                }

                sendVars[keyName] = Variable{
                  Type: "local",
                  Value: ommtypeKey,
                }
                sendVars[valName] = Variable{
                  Type: "local",
                  Value: val,
                }

                interpreted := interpreter(v.ExpAct, cli_params, sendVars, this_vals)

                for k, variable := range interpreted.Variables {
                  vars[k] = variable
                }

                if interpreted.Type == "return" || interpreted.Type == "break" || interpreted.Type == "continue" {
                  return Returner{
                    Type: interpreted.Type,
                    Exp: interpreted.Exp,
                    Variables: vars,
                  }
                }

              }

            case OmmArray:

              for key, val := range it.(OmmArray).Array {

                ommtypeKey := OmmNumber{}
                ommtypeKey.FromGoType(float64(key))

                var sendVars = make(map[string]Variable)

                for k, variable := range vars {
                  sendVars[k] = variable
                }

                sendVars[keyName] = Variable{
                  Type: "local",
                  Value: ommtypeKey,
                }
                sendVars[valName] = Variable{
                  Type: "local",
                  Value: val,
                }

                interpreted := interpreter(v.ExpAct, cli_params, sendVars, this_vals)

                for k, variable := range interpreted.Variables {
                  vars[k] = variable
                }

                if interpreted.Type == "return" || interpreted.Type == "break" || interpreted.Type == "continue" {
                  return Returner{
                    Type: interpreted.Type,
                    Exp: interpreted.Exp,
                    Variables: vars,
                  }
                }

              }

            case OmmString:

              for key, val := range it.(OmmString).ToGoType() {

                ommtypeKey := OmmNumber{}
                ommtypeKey.FromGoType(float64(key))
                ommtypeVal := OmmRune{}
                ommtypeVal.FromGoType(val)

                var sendVars = make(map[string]Variable)

                for k, variable := range vars {
                  sendVars[k] = variable
                }

                sendVars[keyName] = Variable{
                  Type: "local",
                  Value: ommtypeKey,
                }
                sendVars[valName] = Variable{
                  Type: "local",
                  Value: ommtypeVal,
                }

                interpreted := interpreter(v.ExpAct, cli_params, sendVars, this_vals)

                for k, variable := range interpreted.Variables {
                  vars[k] = variable
                }

                if interpreted.Type == "return" || interpreted.Type == "break" || interpreted.Type == "continue" {
                  return Returner{
                    Type: interpreted.Type,
                    Exp: interpreted.Exp,
                    Variables: vars,
                  }
                }

              }

          }

    }
  }

  return Returner{
    Type: "none",
    Variables: vars,
  }
}
