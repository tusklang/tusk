package compiler

import . "lang/types"

func getvars(actions []Action) (map[string][]Action, CompileErr) {

  var vars = make(map[string][]Action)

  for _, v := range actions {
    if v.Type != "var" && v.Type != "declare" { //if it is not an assigner, it must be an error
      return nil, makeCompilerErr("Cannot have anything but a variable declaration outside of a function", v.File, v.Line)
    }
    vars[v.Name] = v.ExpAct
  }

  return vars, nil
}
