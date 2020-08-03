package interpreter

import "fmt"
import "os"
import . "lang/types"

//export OmmPanic
func OmmPanic(err string, line uint64, file string, stacktrace []string) {
  fmt.Println("Panic on line", line, "file", file)
  fmt.Println(err)
  fmt.Println("\nWhen the error was thrown, this was the stack:")
  fmt.Println("  at line", line, "in file", file)
  for i := len(stacktrace) - 1; i >= 0; i-- { //print the stacktrace
    fmt.Println("  " + stacktrace[i])
  }
  os.Exit(1)
}

func (ins *Instance) interpreter(actions []Action, stacktrace []string) Returner {

  var expReturn = false //if it is inside an expression

  if len(actions) == 1 {
    expReturn = true
  }

  for _, v := range actions {
    switch v.Type {

      case "var":

        interpreted := ins.interpreter(v.ExpAct, stacktrace)

        ins.vars[v.Name] = &OmmVar{
          Name: v.Name,
          Value: interpreted.Exp,
        }

        if expReturn {
          variable := *ins.vars[v.Name]
          return Returner{
            Type: "expression",
            Exp: variable.Value,
          }
        }

      case "declare":

        var tmpundef OmmType = undef

        ins.vars[v.Name] = &OmmVar{
          Name: v.Name,
          Value: &tmpundef,
        }

        if expReturn {
          variable := *ins.vars[v.Name]
          return Returner{
            Type: "expression",
            Exp: variable.Value,
          }
        }

      case "delete": //the `delete` statement is inserted at compile-time by the garbage collector
        defer delete(ins.vars, v.Name)

      case "let":

        interpreted := *ins.interpreter(v.ExpAct, stacktrace).Exp

        variable := ins.interpreter(v.First, stacktrace)

        *variable.Exp = interpreted

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: variable.Exp,
          }
        }

      case "log":
        interpreted := ins.interpreter(v.ExpAct, stacktrace)
        fmt.Println((*interpreted.Exp).Format())
      case "print":
        interpreted := ins.interpreter(v.ExpAct, stacktrace)
        fmt.Print((*interpreted.Exp).Format())

      //all of the types
      case "string":   fallthrough
      case "rune":     fallthrough
      case "number":   fallthrough
      case "bool":     fallthrough
      case "function": fallthrough
      case "undef":    fallthrough
      case "c-array":  fallthrough //compile-time calculated array
      case "c-hash":   fallthrough //compile-time calculated hash
      case "thread":

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: &v.Value,
          }
        }

      //arrays, hashes, and protos are a bit different
      case "r-array":

        var nArr = make([]*OmmType, len(v.Array))

        for k, i := range v.Array {
          nArr[k] = ins.interpreter(i, stacktrace).Exp
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

      case "r-hash":

        var nHash = make(map[string]*OmmType)

        for k, i := range v.Hash {
          nHash[k] = ins.interpreter(i, stacktrace).Exp
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
          static[k] = ins.interpreter(i, stacktrace).Exp
        }
        for k, i := range v.Instance {
          instance[k] = ins.interpreter(i, stacktrace).Exp
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

      //////////////////

      case "variable":

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: ins.vars[v.Name].Value,
          }
        }

      case "{": fallthrough
      case "(":

        groupRet := ins.interpreter(v.ExpAct, stacktrace)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: groupRet.Exp,
          }
        }

      case "cast":

        casted := cast(*ins.interpreter(v.ExpAct, stacktrace).Exp, v.Name, stacktrace, v.Line, v.File)

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
      case "!": fallthrough
      case "::": fallthrough
      case "|": fallthrough
      case "&": fallthrough
      case "=>": fallthrough //this is probably not necessary, but i just left it here
      case "<-": fallthrough
      case "<~":

        var firstInterpreted Returner
        var secondInterpreted Returner

        //& and | can be a bit different for bool and bool
        if v.Type == "&" || v.Type == "|" {
          firstInterpreted = ins.interpreter(v.First, stacktrace)
          var computed OmmType

          if (*firstInterpreted.Exp).Type() == "bool" {

            var assumeVal bool

            if v.Type == "&" {
              assumeVal = !isTruthy(*firstInterpreted.Exp)
            } else {
              assumeVal = isTruthy(*firstInterpreted.Exp)
            }

            if assumeVal {
              var retVal = v.Type == "|"
              computed = OmmBool{
                Boolean: &retVal,
              }
            } else {
              secondInterpreted = ins.interpreter(v.Second, stacktrace)

              if (*secondInterpreted.Exp).Type() == "bool" {
                var gobool = isTruthy(*secondInterpreted.Exp)
                computed = OmmBool{
                  Boolean: &gobool,
                }
              } else {
                goto and_or_skip
              }

            }

            if expReturn {
              return Returner{
                Type: "expression",
                Exp: &computed,
              }
            }
          }
        }
        //////////////////////////////////////////////////

        firstInterpreted = ins.interpreter(v.First, stacktrace)
        secondInterpreted = ins.interpreter(v.Second, stacktrace)
        and_or_skip:

        operationFunc, exists := Operations[(*firstInterpreted.Exp).Type() + " " + v.Type + " " + (*secondInterpreted.Exp).Type()]

        if !exists { //if it does not exist, also check the typeof (for protos and objects)
          operationFunc, exists = Operations[(*firstInterpreted.Exp).TypeOf() + " " + v.Type + " " + (*secondInterpreted.Exp).TypeOf()]
        }

        if !exists { //if there is no operation for that type, panic
          OmmPanic("Could not find " + v.Type + " operator for types " + (*firstInterpreted.Exp).TypeOf() + " and " + (*secondInterpreted.Exp).TypeOf(), v.Line, v.File, stacktrace)
        }

        computed := operationFunc(*firstInterpreted.Exp, *secondInterpreted.Exp, ins, stacktrace, v.Line, v.File)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: computed,
          }
        }

      ////////////

      case "await":

        interpreted := ins.interpreter(v.ExpAct, stacktrace).Exp
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
          Exp: ins.interpreter(v.ExpAct, stacktrace).Exp,
        }

      case "condition":

        for _, v := range v.ExpAct {

          truthy := true

          if v.Type == "if" {
            condition := ins.interpreter(v.First, stacktrace)
            truthy = isTruthy(*condition.Exp)
          }

          if truthy {
            interpreted := ins.interpreter(v.ExpAct, stacktrace)

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

        cond := ins.interpreter(v.First, stacktrace)

        for ;isTruthy(*cond.Exp); cond = ins.interpreter(v.First, stacktrace) {

          interpreted := ins.interpreter(v.ExpAct, stacktrace)

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

        it := *ins.interpreter([]Action{ v.First[0] }, stacktrace).Exp
        keyName := v.First[1].Name //get name of key
        valName := v.First[2].Name //get name of val

        switch it.(type) {
          case OmmHash:

            for key, val := range it.(OmmHash).Hash {

              ommtypeKeyn := OmmString{}
              ommtypeKeyn.FromGoType(key)

              var ommtypeKey OmmType = ommtypeKeyn

              ins.vars[keyName] = &OmmVar{
                Name: keyName,
                Value: &ommtypeKey,
              }
              ins.vars[valName] = &OmmVar{
                Name: valName,
                Value: val,
              }

              interpreted := ins.interpreter(v.ExpAct, stacktrace)

              //free the key and val spaces
              delete(ins.vars, keyName)
              delete(ins.vars, valName)
              /////////////////////////////

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

              ins.vars[keyName] = &OmmVar{
                Name: keyName,
                Value: &ommtypeKey,
              }
              ins.vars[valName] = &OmmVar{
                Name: valName,
                Value: val,
              }

              interpreted := ins.interpreter(v.ExpAct, stacktrace)

              //free the key and val spaces
              delete(ins.vars, keyName)
              delete(ins.vars, valName)
              /////////////////////////////

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

          case OmmString:

            for key, val := range it.(OmmString).ToGoType() {

              ommtypeKeyn := OmmNumber{}
              ommtypeKeyn.FromGoType(float64(key))
              ommtypeValr := OmmRune{}
              ommtypeValr.FromGoType(val)

              var ommtypeKey OmmType = ommtypeKeyn
              var ommtypeVal OmmType = ommtypeValr

              ins.vars[keyName] = &OmmVar{
                Name: keyName,
                Value: &ommtypeKey,
              }
              ins.vars[valName] = &OmmVar{
                Name: valName,
                Value: &ommtypeVal,
              }

              interpreted := ins.interpreter(v.ExpAct, stacktrace)

              //free the key and val spaces
              delete(ins.vars, keyName)
              delete(ins.vars, valName)
              /////////////////////////////

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

        }

      case "++":

        variable := ins.interpreter(v.First, stacktrace)

        operationFunc, exists := Operations[(*variable.Exp).Type() + " + number"]

        if !exists { //if there is no operation for that type, panic
          OmmPanic("Could not find + operation for types " + (*variable.Exp).Type() + " and number", v.Line, v.File, stacktrace)
        }

        var onetype OmmType = one
        *variable.Exp = *operationFunc(*variable.Exp, onetype, ins, stacktrace, v.Line, v.File)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: variable.Exp,
          }
        }

      case "--":

        variable := ins.interpreter(v.First, stacktrace)

        operationFunc, exists := Operations[(*variable.Exp).Type() + " - number"]

        if !exists { //if there is no operation for that type, panic
          OmmPanic("Could not find - operation for types " + (*variable.Exp).Type() + " and number", v.Line, v.File, stacktrace)
        }

        var onetype OmmType = one
        *variable.Exp = *operationFunc(*variable.Exp, onetype, ins, stacktrace, v.Line, v.File)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: variable.Exp,
          }
        }

      case "+=": fallthrough
      case "-=": fallthrough
      case "*=": fallthrough
      case "/=": fallthrough
      case "%=": fallthrough
      case "^=":

        variable := ins.interpreter(v.First, stacktrace)
        interpreted := *ins.interpreter(v.Second, stacktrace).Exp

        operationFunc, exists := Operations[(*variable.Exp).Type() + " " + string(v.Type[0]) + " " + interpreted.Type()]

        if !exists { //if there is no operation for that type, panic
          OmmPanic("Could not find " + string(v.Type[0]) + " operation for types " + (*variable.Exp).Type() + " and " + interpreted.Type(), v.Line, v.File, stacktrace)
        }

        *variable.Exp = *operationFunc(*variable.Exp, interpreted, ins, stacktrace, v.Line, v.File)

        if expReturn {
          return Returner{
            Type: "expression",
            Exp: variable.Exp,
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
