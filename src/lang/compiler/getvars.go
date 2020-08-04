package compiler

import . "lang/types"

func getvars(actions []Action) (map[string][]Action, CompileErr) {

  var vars = make(map[string][]Action)

  for _, v := range actions {
    if v.Type != "var" && v.Type != "declare" && v.Type != "ovld" { //if it is not an assigner or overloader, it must be an error
      return nil, makeCompilerErr("Cannot have anything but a variable declaration or overloader outside of a function", v.File, v.Line)
    }

    if v.Type == "ovld" {
      if _, exists := vars[v.Name]; !exists { //if it does not exist yet, declare undefined yet
        return nil, makeCompilerErr("Undefined overloader: " + v.Name[1:], v.File, v.Line)
      }
      vars["ovld/" + v.Name] = v.ExpAct //set it to an overloader
      continue
    }

    if _, exists := vars[v.Name]; exists { //if the given global name already exists, throw an error
      return nil, makeCompilerErr("Duplicate global name was detected", v.File, v.Line)
    }

    vars[v.Name] = v.ExpAct
  }

  return vars, nil
}
