package interpreter

import "fmt"
import "os"
import . "lang/types"

//export OmmPanic
func OmmPanic(err string, line uint64, file string, stacktrace []string) {
  fmt.Println("Panic on line", line, "file", file)
  fmt.Println(err)
  fmt.Println("\nWhen the error was thrown, this was the stack:")
  for i := len(stacktrace) - 1; i >= 0; i-- { //print the stacktrace
    fmt.Println("  " + stacktrace[i])
  }
  os.Exit(1)
}

func interpreter(actions []Action, cli_params CliParams, stacktrace []string) Returner {

  var expReturn = false //if it is inside an expression

  if len(actions) == 1 {
    expReturn = true
  }

  for _, v := range actions {
    switch v.Type {

      case "var":

        interpreted := interpreter(v.ExpAct, cli_params, stacktrace)

        vars[v.Name] = Variable{
          Type: "var",
          Value: interpreted.Exp,
        }

        if expReturn {
          variable := vars[v.Name]
          return Returner{
            Type: "expression",
            Exp: variable.Value,
          }
        }

      case "declare":

        var tmpundef OmmType = undef

        vars[v.Name] = Variable{
          Type: "var",
          Value: &tmpundef,
        }

        if expReturn {
          variable := vars[v.Name]
          return Returner{
            Type: "expression",
            Exp: variable.Value,
          }
        }

      case "let":

        interpreted := *interpreter(v.ExpAct, cli_params, stacktrace).Exp

        variable := interpreter(v.First, cli_params, stacktrace)

        *variable.Exp = interpreted

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: variable.Exp,
          }
        }

      case "log":
        interpreted := interpreter(v.ExpAct, cli_params, stacktrace)
        fmt.Println((*interpreted.Exp).Format())
      case "print":
        interpreted := interpreter(v.ExpAct, cli_params, stacktrace)
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

      //arrays, hashes, and structures are a bit different
      case "array":

        var nArr = make([]*OmmType, len(v.Array))

        for k, i := range v.Array {
          nArr[k] = interpreter(i, cli_params, stacktrace).Exp
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
          nHash[k] = interpreter(i, cli_params, stacktrace).Exp
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

      case "proto":

        var static = make(map[string]*OmmType)
        var instance = make(map[string]*OmmType)

        for k, i := range v.Static {
          static[k] = interpreter(i, cli_params, stacktrace).Exp
        }
        for k, i := range v.Instance {
          instance[k] = interpreter(i, cli_params, stacktrace).Exp
        }

        var proto = OmmProto{
          ProtoName: v.Name,
        }
        proto.Set(static, instance)

        if expReturn {
          var ommtype OmmType = proto
          return Returner{
            Type: "expression",
            Exp: &ommtype,
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

        groupRet := interpreter(v.ExpAct, cli_params, stacktrace)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: groupRet.Exp,
          }
        }

      case "cast":

        casted := cast(*interpreter(v.ExpAct, cli_params, stacktrace).Exp, v.Name, stacktrace, v.Line, v.File)

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
      case "<-": fallthrough
      case "<~":

        firstInterpreted := interpreter(v.First, cli_params, stacktrace)
        secondInterpreted := interpreter(v.Second, cli_params, stacktrace)

        operationFunc, exists := operations[(*firstInterpreted.Exp).Type() + " " + v.Type + " " + (*secondInterpreted.Exp).Type()]

        if !exists { //if there is no operation for that type, panic
          OmmPanic("Could not find " + v.Type + " operation for types " + (*firstInterpreted.Exp).Type() + " and " + (*secondInterpreted.Exp).Type(), v.Line, v.File, stacktrace)
        }

        computed := operationFunc(*firstInterpreted.Exp, *secondInterpreted.Exp, cli_params, stacktrace, v.Line, v.File)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: computed,
          }
        }

      ////////////

      case "await":

        interpreted := interpreter(v.ExpAct, cli_params, stacktrace).Exp
        var awaited OmmType

        switch (*interpreted).(type) {
          case OmmThread:

            //put the new value back into the given interpreted pointer
            thread := (*interpreted).(OmmThread)
            thread.WaitFor()
            *interpreted = thread
            ///////////////////////////////////////////////////////////

            awaited = *thread.Returned.Exp
          default:
            awaited = *interpreted
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
          Exp: interpreter(v.ExpAct, cli_params, stacktrace).Exp,
        }

      case "condition":

        for _, v := range v.ExpAct {

          truthy := true

          if v.Type == "if" {
            condition := interpreter(v.First, cli_params, stacktrace)
            truthy = isTruthy(*condition.Exp)
          }

          if truthy {
            interpreted := interpreter(v.ExpAct, cli_params, stacktrace)

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

        cond := interpreter(v.First, cli_params, stacktrace)

        for ;isTruthy(*cond.Exp); cond = interpreter(v.First, cli_params, stacktrace) {

          interpreted := interpreter(v.ExpAct, cli_params, stacktrace)

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

        it := *interpreter([]Action{ v.First[0] }, cli_params, stacktrace).Exp
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

              interpreted := interpreter(v.ExpAct, cli_params, stacktrace)

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

              interpreted := interpreter(v.ExpAct, cli_params, stacktrace)

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

              interpreted := interpreter(v.ExpAct, cli_params, stacktrace)

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
