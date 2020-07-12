package compiler

import . "lang/types"

func getvars(actions []Action, dir string) map[string][]Action {

  var vars = make(map[string][]Action)

  for _, v := range actions {
    if v.Type != "var" { //if it is not an assigner, it must be an error
      compilerErr("Cannot have anything but a variable declaration outside of a function", dir, v.Line)
    }
    vars[v.Name] = v.ExpAct
  }

  return vars
}
