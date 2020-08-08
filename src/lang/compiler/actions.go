package compiler

import "runtime"
import . "lang/types"
import . "lang/interpreter"

func actionizer(operations []Operation) ([]Action, CompileErr) {

  var actions []Action
  var e CompileErr

  for _, v := range operations {

    var left []Action
    var right []Action

    if v.Left != nil {
      left, e = actionizer([]Operation{ *v.Left })

      if e != nil {
        return []Action{}, e
      }
    }
    if v.Right != nil {
      right, e = actionizer([]Operation{ *v.Right })

      if e != nil {
        return []Action{}, e
      }
    }

    switch v.Type {
      case "~":

        var statements = []string{ "var", "log", "print", "if", "elif", "else", "while", "each", "include", "function", "return", "await", "proto", "static", "instance", "ifwin", "ifnwin", "ovld" } //list of statements

        var hasStatement bool = false

        for _, val := range statements {
          if val == (*v.Left).Item.Token.Name {

            switch val {

              case "include":
                if right[0].Type != "string" {
                  return []Action{}, makeCompilerErr("Expected a string after \"include\"", v.File, v.Line)
                }

                includeFiles, e := includer(right[0].Value.(OmmString).ToGoType(), v.Line, v.File)

                if e != nil {
                  return []Action{}, e
                }

                for _, acts := range includeFiles {
                  actions = append(actions, acts...)
                }

              case "function":
                if right[0].Type != "=>" {
                  return []Action{}, makeCompilerErr("Functions need a parameter list and a function body", v.File, right[0].Line)
                }

                var typeList []string
                var paramList []string

                for _, p := range right[0].First[0].ExpAct {
                  if p.Type != "cast" || p.ExpAct[0].Type != "variable" {
                    return []Action{}, makeCompilerErr("Function parameter lists can only have typed variables", v.File, right[0].Line)
                  }

                  typeList = append(typeList, p.Name)
                  paramList = append(paramList, p.ExpAct[0].Name)
                }

                actions = append(actions, Action{
                  Type: "function",
                  Value: OmmFunc{
                    Overloads: []Overload{
                      Overload{
                        Params: paramList,
                        Types: typeList,
                        Body: right[0].Second,
                      },
                    },
                  },
                  File: v.File,
                  Line: v.Line,
                })

              case "if":

                if right[0].Type != "=>" {
                  return []Action{}, makeCompilerErr("If statements need a condition and a body", v.File, right[0].Line)
                }

                actions = append(actions, Action{
                  Type: "condition",
                  ExpAct: []Action{ Action{
                    Type: "if",
                    First: right[0].First,
                    ExpAct: right[0].Second,
                  } },
                  File: v.File,
                  Line: v.Line,
                })

              case "elif":

                if right[0].Type != "=>" {
                  return []Action{}, makeCompilerErr("Elif statements need a condition and a body", v.File, right[0].Line)
                }

                if len(actions) == 0 || actions[len(actions) - 1].Type != "condition" {
                  return []Action{}, makeCompilerErr("Unexpected elif statement", v.File, right[0].Line)
                }

                //append to the previous conditional statement
                actions[len(actions) - 1].ExpAct = append(actions[len(actions) - 1].ExpAct, Action{
                  Type: "if",
                  First: right[0].First,
                  ExpAct: right[0].Second,
                })

              case "else":

                if len(actions) == 0 || actions[len(actions) - 1].Type != "condition" {
                  return []Action{}, makeCompilerErr("Unexpected else statement", v.File, right[0].Line)
                }

                //append to the previous conditional statement
                actions[len(actions) - 1].ExpAct = append(actions[len(actions) - 1].ExpAct, Action{
                  Type: "else",
                  ExpAct: right,
                })

              case "while":
                if right[0].Type != "=>" {
                  return []Action{}, makeCompilerErr("While loops need a condition and a body", v.File, right[0].Line)
                }

                actions = append(actions, Action{
                  Type: val,
                  First: right[0].First,
                  ExpAct: right[0].Second,
                  File: v.File,
                  Line: v.Line,
                })

              case "each":
                if right[0].Type != "=>" {
                  return []Action{}, makeCompilerErr("Each loops need a condition and a body", v.File, right[0].Line)
                }

                if len(right[0].First[0].ExpAct) != 3 {
                  return []Action{}, makeCompilerErr("Each loops must look like this: each(iterator, key, value)", v.File, right[0].Line)
                }

                for _, n := range right[0].First[0].ExpAct[1:] {
                  if n.Type != "variable" {
                    return []Action{}, makeCompilerErr("Key or value was not given as a variable", v.File, right[0].Line)
                  }
                }

                actions = append(actions, Action{
                  Type: val,
                  First: right[0].First[0].ExpAct, //because it doesnt matter if they use a { or (
                  ExpAct: right[0].Second,
                  File: v.File,
                  Line: v.Line,
                })

              case "var":

                if right[0].Type == "variable" { //the dev is declaring is like "var a" (meaning declare a)
                  actions = append(actions, Action{
                    Type: "declare",
                    Name: right[0].Name,
                    File: v.File,
                    Line: v.Line,
                  })
                } else {
                  if right[0].Type != "let" {
                    return []Action{}, makeCompilerErr("Expected a assigner statement after var", v.File, right[0].Line)
                  }

                  if right[0].First[0].Type != "variable" {
                    return []Action{}, makeCompilerErr("Cannot use :: operator in variable declaration", v.File, right[0].Line)
                  }

                  actions = append(actions, Action{
                    Type: val,
                    Name: right[0].First[0].Name,
                    ExpAct: right[0].ExpAct,
                    File: v.File,
                    Line: v.Line,
                  })
                }

              case "proto":

                if len(right) == 0 {
                  return []Action{}, makeCompilerErr("Prototypes require a body", v.File, right[0].Line)
                }

                if right[0].Type != "{" {
                  return []Action{}, makeCompilerErr("Prototype bodies can only be curly brace enclosed", v.File, right[0].Line)
                }

                var (
                  static = make(map[string]*OmmType)
                  instance = make(map[string]*OmmType)
                )
                var body = right[0].ExpAct //get the struct body

                for i := range body {

                  if body[i].Type != "static" && body[i].Type != "instance" { //if it does not name static or instance, automatically make it instance
                    body[i] = Action{
                      Type: "instance",
                      ExpAct: []Action{ body[i] },
                      File: body[i].File,
                      Line: body[i].Line,
                    }
                  }

                  name := body[i].ExpAct[0].Name

                  if body[i].ExpAct[0].Type == "var" {

                    if len(body[i].ExpAct[0].ExpAct) == 0 || body[i].ExpAct[0].ExpAct[0].Value == nil {
                      return []Action{}, makeCompilerErr("Cannot have compound types at the golabl scope of a prototype", v.File, right[0].Line)
                    }

                    if body[i].Type == "static" {
                      static[name] = &body[i].ExpAct[0].ExpAct[0].Value
                    } else {
                      instance[name] = &body[i].ExpAct[0].ExpAct[0].Value
                    }

                  } else if body[i].ExpAct[0].Type == "declare" {

                    var tmp OmmType = OmmUndef{}

                    if body[i].Type == "static" {
                      static[name] = &tmp
                    } else {
                      instance[name] = &tmp
                    }
                  } else {
                    return []Action{}, makeCompilerErr("Prototype bodies can only have variable assignments and declarations", v.File, right[0].Line)
                  }
                }

                actions = append(actions, Action{
                  Type: "proto",
                  Value: OmmProto{
                    Static: static,
                    Instance: instance,
                  },
                  File: v.File,
                  Line: v.Line,
                })

              case "ifwin":

                if runtime.GOOS == "windows" {
                  actions = append(actions, right...)
                }

              case "ifnwin":

                if runtime.GOOS != "windows" {
                  actions = append(actions, right...)
                }

              case "ovld":

                if right[0].Type != "let" {
                  return []Action{}, makeCompilerErr("Expected a assigner statement after ovld", v.File, right[0].Line)
                }

                if right[0].First[0].Type != "variable" {
                  return []Action{}, makeCompilerErr("Cannot use :: operator in an overloader", v.File, right[0].Line)
                }

                if right[0].ExpAct[0].Type != "function" {
                  return []Action{}, makeCompilerErr("Cannot overload a " + right[0].ExpAct[0].Type, v.File, right[0].Line)
                }

                for _, p := range right[0].ExpAct[0].Value.(OmmFunc).Overloads[0].Types {
                  if p == "any" {
                    return []Action{}, makeCompilerErr("Cannot use \"any\" type in a overload (may cause bugs)", v.File, right[0].Line)
                  }
                }

                right[0].Type = "ovld"

                actions = append(actions, Action{
                  Type: val,
                  Name: right[0].First[0].Name,
                  ExpAct: right[0].ExpAct,
                  File: v.File,
                  Line: v.Line,
                })

              default:

                actions = append(actions, Action{
                  Type: val,
                  ExpAct: right,
                  File: v.File,
                  Line: v.Line,
                })

            }

            hasStatement = true
          }
        }

        if !hasStatement {
          return []Action{}, makeCompilerErr("\"" + (*v.Left).Item.Token.Name + "\" is not a statement", v.File, v.Line)
        }

      case ":=":

        if left[0].Type != "variable" {
          return []Action{}, makeCompilerErr("Expected a variable statement before := operator", v.File, right[0].Line)
        }

        actions = append(actions, Action{
          Type: "var",
          Name: left[0].Name,
          ExpAct: right,
          File: v.File,
          Line: v.Line,
        })

      case "=":

        actions = append(actions, Action{
          Type: "let",
          First: left,
          ExpAct: right,
          File: v.File,
          Line: v.Line,
        })

      case "->":

        castType := v.Left.Item.Token.Name[1:]

        actions = append(actions, Action{
          Type: "cast",
          Name: castType, //type to cast into
          ExpAct: right,
          File: v.File,
          Line: v.Line,
        })

      //all of these operations have the same way of appending
      case "::":

        //if it is ::, and the next action is a variable, then convert to a string
        //to get index of a variable's value, use ::()
        //for example,
        //  var a: [:
        //    "hello": "world",
        //  :]
        //  log a::hello ; would log "world"
        //
        //  var idx: "hello"
        //  log a::(idx) ; would log "world" as well
        //  log a::idx ; //would cause a panic error

        if len(right) == 0 { //safeguard
          return []Action{}, makeCompilerErr("No value was found right of a :: operator", v.File, v.Line)
        }


        if right[0].Type == "variable" {
          var str = OmmString{}
          str.FromGoType(right[0].Name[1:]) //remove the $ from the varname
          right[0] = Action{
            Type: "string",
            Value: str,
            File: right[0].File,
            Line: right[0].Line,
          }
        }

        fallthrough
      case "+": fallthrough
      case "-": fallthrough
      case "*": fallthrough
      case "/": fallthrough
      case "%": fallthrough
      case "^": fallthrough
      case "==": fallthrough
      case "!=": fallthrough
      case ">": fallthrough
      case "<": fallthrough
      case ">=": fallthrough
      case "<=": fallthrough
      case "!": fallthrough
      case "&": fallthrough
      case "|": fallthrough
      case "=>": fallthrough
      case "<-": fallthrough
      case "<~":

        actions = append(actions, Action{
          Type: v.Type,
          First: left,
          Second: right,
          File: v.File,
          Line: v.Line,
        })
      ////////////////////////////////////////////////////////

      case "++": fallthrough
      case "--":

        if len(left) == 0 || (left[0].Type != "variable" && left[0].Type != "::") {
          return []Action{}, makeCompilerErr("Must have a variable before an increment or decrement", v.File, v.Line)
        }

        actions = append(actions, Action{
          Type: v.Type,
          First: left,
          File: v.File,
          Line: v.Line,
        })

      case "+=": fallthrough
      case "-=": fallthrough
      case "*=": fallthrough
      case "/=": fallthrough
      case "%=": fallthrough
      case "^=":

        if len(left) == 0 || (left[0].Type != "variable" && left[0].Type != "::") {
          return []Action{}, makeCompilerErr("Must have a variable before an assignment operator", v.File, v.Line)
        }
        if len(right) == 0 {
          return []Action{}, makeCompilerErr("Could not find a value after " + v.Type, v.File, v.Line)
        }

        actions = append(actions, Action{
          Type: v.Type,
          First: left,
          Second: right,
          File: v.File,
          Line: v.Line,
        })

      case "break": fallthrough
      case "continue":

        actions = append(actions, Action{
          Type: v.Type,
          File: v.File,
          Line: v.Line,
        })

      case "none":
        vActs, e := valueActions(v.Item)

        if e != nil {
          return []Action{}, e
        }

        actions = append(actions, vActs)
    }
  }

  return actions, nil
}
