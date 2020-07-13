package interpreter

import "fmt"
import "os"
import . "lang/types"

func ommPanic(err string, line uint64, file string) {
  fmt.Println("Panic on line", line, "file", file)
  fmt.Println(err)
  os.Exit(1)
}

func interpreter(actions []Action, cli_params CliParams) Returner {

  var expReturn = false //if it is inside an expression

  if len(actions) == 1 {
    expReturn = true
  }

  for _, v := range actions {
    switch v.Type {

      case "var":

        interpreted := interpreter(v.ExpAct, cli_params)

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

        interpreted := *interpreter(v.ExpAct, cli_params).Exp

        variable := interpreter(v.First, cli_params)

        *variable.Exp = interpreted

      case "log":
        interpreted := interpreter(v.ExpAct, cli_params)
        fmt.Println((*interpreted.Exp).Format())
      case "print":
        interpreted := interpreter(v.ExpAct, cli_params)
        fmt.Print((*interpreted.Exp).Format())

      //all of the types
      case "string": fallthrough
      case "rune": fallthrough
      case "number": fallthrough
      case "bool": fallthrough
      case "function": fallthrough
      case "undef": fallthrough
      case "thread":

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: &v.Value,
          }
        }

      //arrays and hashes are a bit different
      case "array":

        var nArr = make([]*OmmType, len(v.Array))

        for k, i := range v.Array {
          nArr[k] = interpreter(i, cli_params).Exp
        }

        var ommType OmmType = OmmArray{
          Array: nArr,
          Length: uint64(len(v.Array)),
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: &ommType,
          }
        }

      case "hash":

        var nHash = make(map[string]*OmmType)

        for k, i := range v.Hash {
          nHash[k] = interpreter(i, cli_params).Exp
        }

        var ommType OmmType = OmmHash{
          Hash: nHash,
          Length: uint64(len(v.Hash)),
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: &ommType,
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

      case "{": fallthrough
      case "(":

        groupRet := interpreter(v.ExpAct, cli_params)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: groupRet.Exp,
          }
        }

      case "->":

        casted := cast(*interpreter(v.ExpAct, cli_params).Exp, v.Name)

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

        firstInterpreted := interpreter(v.First, cli_params)
        secondInterpreted := interpreter(v.Second, cli_params)

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

      case "await":

        interpreted := *interpreter(v.ExpAct, cli_params).Exp
        var awaited OmmType

        switch interpreted.(type) {
          case OmmThread:
            awaited = *interpreted.(OmmThread).WaitFor().Exp
          default:
            awaited = interpreted
        }

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: &awaited,
          }
        }

      case "break": fallthrough
      case "continue":

        return Returner{
          Type: v.Type,
        }

      case "return":

        return Returner{
          Type: "return",
          Exp: interpreter(v.ExpAct, cli_params).Exp,
        }

      case "condition":

        for _, v := range v.ExpAct {

          truthy := true

          if v.Type == "if" {
            condition := interpreter(v.First, cli_params)
            truthy = isTruthy(*condition.Exp)
          }

          if truthy {
            interpreted := interpreter(v.ExpAct, cli_params)

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

        cond := interpreter(v.First, cli_params)

        for ;isTruthy(*cond.Exp); cond = interpreter(v.First, cli_params) {

          interpreted := interpreter(v.ExpAct, cli_params)

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

        it := *interpreter([]Action{ v.First[0] }, cli_params).Exp
        keyName := v.First[1].Name //get name of key
        valName := v.First[2].Name //get name of val

        switch it.(type) {
          case OmmHash:

            for key, val := range it.(OmmHash).Hash {

              ommtypeKeyn := OmmString{}
              ommtypeKeyn.FromGoType(key)

              var ommtypeKey OmmType = ommtypeKeyn

              vars[keyName] = Variable{
                Type: "local",
                Value: &ommtypeKey,
              }
              vars[valName] = Variable{
                Type: "local",
                Value: val,
              }

              interpreted := interpreter(v.ExpAct, cli_params)

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

              vars[keyName] = Variable{
                Type: "local",
                Value: &ommtypeKey,
              }
              vars[valName] = Variable{
                Type: "local",
                Value: val,
              }

              interpreted := interpreter(v.ExpAct, cli_params)

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

              vars[keyName] = Variable{
                Type: "local",
                Value: &ommtypeKey,
              }
              vars[valName] = Variable{
                Type: "local",
                Value: &ommtypeVal,
              }

              interpreted := interpreter(v.ExpAct, cli_params)

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
