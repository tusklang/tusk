package compiler

import "fmt"
import "os"
import "io/ioutil"
import "strings"

import "lang/interpreter"
import . "lang/types"

type _CompileErr struct {
  Msg   string
  FName string
  Line  uint64
}

func (e _CompileErr) Print() {
  fmt.Println("Error while compiling", e.FName, "on line", e.Line)
  fmt.Println(e.Msg)
  os.Exit(1)
}

type CompileErr interface {
  Print()
}

func makeCompilerErr(msg, fname string, line uint64) CompileErr {
  return _CompileErr{
    Msg: msg,
    FName: fname,
    Line: line,
  }
}

func inclCompile(file, filename string, compileall bool) ([]Action, CompileErr) {

  lex, e := lexer(file, filename)

  if e != nil {
    return []Action{}, e
  }

  if compileall { //if the dev wants to compile the entire directory

    files, _ := ioutil.ReadDir(".")

    for _, v := range files {
      if strings.HasSuffix(v.Name(), ".omm") {
        lex = append([]Lex{
          Lex{
            Name: "include",
            Type: "keyword",
          },
          Lex{
            Name: "~",
            Type: "operation",
          },
          Lex{
            Name: "\"" + v.Name() + "\"",
            Type: "expression value",
          },
          Lex{
            Name: "$term",
            Type: "?none",
          },
        }, lex...)
      }
    }

  }

  groups, e := makeGroups(lex)

  if e != nil {
    return []Action{}, e
  }

  operations, e := makeOperations(groups)

  if e != nil {
    return []Action{}, e
  }

  actions, e := actionizer(operations)

  return actions, e
}

//export Compile
func Compile(file, filename string, compileall, isoat bool) ([]Action, map[string][]Action, CompileErr) {

  var e CompileErr

  actions, e := inclCompile(file, filename, compileall)

  if e != nil {
    return nil, nil, e
  }

  if e != nil {
    return []Action{}, nil, e
  }

  //a bunch of validations and initializers
  e = has_non_global_prototypes(actions, true)
  if e != nil {
    return []Action{}, nil, e
  }
  put_proto_types(actions)
  e = validate_types(actions)
  if e != nil {
    return []Action{}, nil, e
  }
  /////////////////////////////////////////

  vars, e := getvars(actions)
  if e != nil {
    return nil, nil, e
  }

  //make each var have only it's name
  var varnames = make(map[string]string)

  varnames["$__dirname"] = "$__dirname"
  varnames["$argv"] = "$argv"

  for k := range vars {

    if len(vars[k]) == 0 { //skip for declares
      varnames[k] = k
      continue
    }

    //ensure that the globals do not have any compound types (such as operations)
    if vars[k][0].Value == nil {
      return nil, nil, makeCompilerErr("Cannot have compound types at the global scope", vars[k][0].File, vars[k][0].Line)
    }

    varnames[k] = k
  }

  //also account for the gofuncs
  for k := range interpreter.GoFuncs {
    varnames["$" + k] = "$" + k
  }

  for k := range vars {
    _, e = changevarnames(vars[k], varnames)
    if e != nil {
      return []Action{}, nil, e
    }

    vars[k] = insert_garbage_collectors(vars[k])
  }

  return actions, vars, nil
}
