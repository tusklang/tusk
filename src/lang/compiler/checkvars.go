package compiler

import . "lang/types"

func checkvars(actions []Action, dir string, vars map[string]string) {

  var curVars = map[string]string{}

  for k, v := range vars {
    curVars[k] = v
  }

  for _, v := range actions {
    if v.Type == "let" || v.Type == "global" || v.Type == "local" {
      curVars[v.Name] = v.Type
    }

    if v.Type == "variable" {
      if _, exists := curVars[v.Name]; !exists {
        compilerErr(v.Name[1:] /* remove the $ from the variable name */ + " was not declared", dir, v.Line)
      }
    }

    if v.Type == "function" {
      for _, p := range v.Value.(OmmFunc).Params { //add the params to the current variables
        curVars[p.Name] = "local"
      }
      checkvars(v.Value.(OmmFunc).Body, dir, curVars)
    }

    //perform checkvars on all of the sub actions
    checkvars(v.ExpAct, dir, curVars)
    checkvars(v.First, dir, curVars)
    checkvars(v.Second, dir, curVars)
    checkvars(v.Degree, dir, curVars)
    /////////////////////////////////////////////

  }

  //send the globals back to the outer scope
  for k, v := range curVars {
    if v == "global" {
      vars[k] = v
    }
  }
}
