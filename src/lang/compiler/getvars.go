package compiler

import . "lang/types"

func getvars(actions []Action, dir string) map[string][]Action {

  var vars = make(map[string][]Action)

  for _, v := range actions {
    if v.Type != "local" && v.Type != "global" { //if it is not an assigner, it must be an error
      compilerErr("Cannot have anything but an assigner outside of a function", dir, v.Line)
    }
    vars[v.Name] = v.ExpAct
  }

  return vars
}
