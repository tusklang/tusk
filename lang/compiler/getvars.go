package compiler

import . "github.com/omm-lang/omm/lang/types"

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

      //since the variables are passed as a map, each overload name must be different
      //omm inserts spaces after each overload to differentiate them
      var space_amt = ""

      _, exists := vars["ovld/" + v.Name]
      for ;exists; _, exists = vars["ovld/" + v.Name + space_amt] {
        space_amt+=" "
      }

      vars["ovld/" + v.Name + space_amt] = v.ExpAct //set it to an overloader
      continue
    }

    vars[v.Name] = v.ExpAct
  }

  return vars, nil
}
