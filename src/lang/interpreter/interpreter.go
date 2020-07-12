package interpreter

import "fmt"
import "os"
import . "lang/types"

func ommPanic(err string, line uint64, file string) {
  fmt.Println("Panic on line", line, "file", file)
  fmt.Println(err)
  os.Exit(1)
}

func interpreter(actions []Action, cli_params CliParams, vars map[string]Variable, this_vals []Action) Returner {

  var expReturn = false //if it is inside an expression

  if len(actions) == 1 {
    expReturn = true
  }

  for _, v := range actions {
    switch v.Type {

      case "var":

        interpreted := interpreter(v.ExpAct, cli_params, vars, this_vals)

        vars[v.Name] = Variable{
          Type: v.Type,
          Value: interpreted.Exp,
        }

        if expReturn {
          variable := vars[v.Name]
          return Returner{
            Type: "expression",
            Exp: variable.Value,
          }
        }

      case "let":

        interpreted := *interpreter(v.ExpAct, cli_params, vars, this_vals).Exp

        variable := interpreter(v.First, cli_params, vars, this_vals)

        *variable.Exp = interpreted

      case "log":
        interpreted := interpreter(v.ExpAct, cli_params, vars, this_vals)
        fmt.Println((*interpreted.Exp).Format())
      case "print":
        interpreted := interpreter(v.ExpAct, cli_params, vars, this_vals)
        fmt.Print((*interpreted.Exp).Format())

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
            Exp: &v.Value,
          }
        }
      //////////////////

      case "variable":

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: vars[v.Name].Value,
          }
        }

      case "{":

        var passedVars = make(map[string]Variable)

        for k, variable := range vars {
          passedVars[k] = variable
        }

        groupRet := interpreter(v.ExpAct, cli_params, passedVars, this_vals)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: groupRet.Exp,
          }
        }

      case "(":

        groupRet := interpreter(v.ExpAct, cli_params, vars, this_vals)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: groupRet.Exp,
          }
        }

      case "->":

        casted := cast(*interpreter(v.ExpAct, cli_params, vars, this_vals).Exp, v.Name)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: casted,
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
      case "=>": fallthrough //this is probably not necessary, but i just left it here
      case "sync": fallthrough
      case "async":
        firstInterpreted := interpreter(v.First, cli_params, vars, this_vals)
        secondInterpreted := interpreter(v.Second, cli_params, vars, this_vals)

        operationFunc, exists := operations[(*firstInterpreted.Exp).Type() + " " + v.Type + " " + (*secondInterpreted.Exp).Type()]

        if !exists { //if there is no operation for that type, panic
          ommPanic("Could not find " + v.Type + " operation for types " + (*firstInterpreted.Exp).Type() + " and " + (*secondInterpreted.Exp).Type(), v.Line, v.File)
        }

        computed := operationFunc(*firstInterpreted.Exp, *secondInterpreted.Exp, cli_params, v.Line, v.File)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: computed,
          }
        }

      ////////////

      case "break": fallthrough
      case "continue":

        return Returner{
          Type: v.Type,
        }

      case "return":

        return Returner{
          Type: "return",
          Exp: interpreter(v.ExpAct, cli_params, vars, this_vals).Exp,
        }

      case "condition":

        for _, v := range v.ExpAct {

          truthy := true

          if v.Type == "if" {
            condition := interpreter(v.First, cli_params, vars, this_vals)
            truthy = isTruthy(*condition.Exp)
          }

          if truthy {
            interpreted := interpreter(v.ExpAct, cli_params, vars, this_vals)

            if interpreted.Type == "return" || interpreted.Type == "break" || interpreted.Type == "continue" {
              return Returner{
                Type: interpreted.Type,
                Exp: interpreted.Exp,
              }
            }

            break
          }
        }

      case "while":

        cond := interpreter(v.First, cli_params, vars, this_vals)

        for ;isTruthy(*cond.Exp); cond = interpreter(v.First, cli_params, vars, this_vals) {

          interpreted := interpreter(v.ExpAct, cli_params, vars, this_vals)

          if interpreted.Type == "return" {
            return Returner{
              Type: interpreted.Type,
              Exp: interpreted.Exp,
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

        it := *interpreter([]Action{ v.First[0] }, cli_params, vars, this_vals).Exp
        keyName := v.First[1].Name //get name of key
        valName := v.First[2].Name //get name of val

        switch it.(type) {
          case OmmHash:

            for key, val := range it.(OmmHash).Hash {

              ommtypeKeyn := OmmString{}
              ommtypeKeyn.FromGoType(key)

              var ommtypeKey OmmType = ommtypeKeyn

              var sendVars = make(map[string]Variable)

              for k, variable := range vars {
                sendVars[k] = variable
              }

              sendVars[keyName] = Variable{
                Type: "local",
                Value: &ommtypeKey,
              }
              sendVars[valName] = Variable{
                Type: "local",
                Value: val,
              }

              interpreted := interpreter(v.ExpAct, cli_params, sendVars, this_vals)

              if interpreted.Type == "return" {
                return Returner{
                  Type: interpreted.Type,
                  Exp: interpreted.Exp,
                }
              }

              if interpreted.Type == "break" {
                break
              }
              if interpreted.Type == "continue" {
                continue;
              }

            }

          case OmmArray:

            for key, val := range it.(OmmArray).Array {

              ommtypeKeyn := OmmNumber{}
              ommtypeKeyn.FromGoType(float64(key))

              var ommtypeKey OmmType = ommtypeKeyn

              var sendVars = make(map[string]Variable)

              for k, variable := range vars {
                sendVars[k] = variable
              }

              sendVars[keyName] = Variable{
                Type: "local",
                Value: &ommtypeKey,
              }
              sendVars[valName] = Variable{
                Type: "local",
                Value: val,
              }

              interpreted := interpreter(v.ExpAct, cli_params, sendVars, this_vals)

              if interpreted.Type == "return" || interpreted.Type == "break" || interpreted.Type == "continue" {
                return Returner{
                  Type: interpreted.Type,
                  Exp: interpreted.Exp,
                }
              }

            }

          case OmmString:

            for key, val := range it.(OmmString).ToGoType() {

              ommtypeKeyn := OmmNumber{}
              ommtypeKeyn.FromGoType(float64(key))
              ommtypeValr := OmmRune{}
              ommtypeValr.FromGoType(val)

              var ommtypeKey OmmType = ommtypeKeyn
              var ommtypeVal OmmType = ommtypeValr

              var sendVars = make(map[string]Variable)

              for k, variable := range vars {
                sendVars[k] = variable
              }

              sendVars[keyName] = Variable{
                Type: "local",
                Value: &ommtypeKey,
              }
              sendVars[valName] = Variable{
                Type: "local",
                Value: &ommtypeVal,
              }

              interpreted := interpreter(v.ExpAct, cli_params, sendVars, this_vals)

              if interpreted.Type == "return" || interpreted.Type == "break" || interpreted.Type == "continue" {
                return Returner{
                  Type: interpreted.Type,
                  Exp: interpreted.Exp,
                }
              }

            }

        }

    }
  }

  var undefval OmmType = OmmUndef{}

  return Returner{
    Type: "none",
    Exp: &undefval,
  }
}
