package compiler

import . "lang/types"

var types = []string{ "string", "rune", "number", "bool", "hash", "array" }

func validate_types(actions []Action) {

  //function to make sure the typecasts (and fargc) work

  for _, v := range actions {

    if v.Type == "cast" {
      for _, t := range types {
        if t == v.Name { //if the type exists, do not throw an error
          goto cast_noErr
        }
      }
      compilerErr("\"" + v.Name + "\" is not a type", v.File, v.Line)
      cast_noErr:
    }

    if v.Type == "function" {
      validate_types(v.Value.(OmmFunc).Body)
      continue
    }
    if v.Type == "proto" {

      for i := range v.Static {
        validate_types(v.Static[i])
      }
      for i := range v.Instance {
        validate_types(v.Instance[i])
      }

      continue
    }

    //perform checkvars on all of the sub actions
    validate_types(v.ExpAct)
    validate_types(v.First)
    validate_types(v.Second)
    validate_types(v.Degree)

    //also do it for the arrays and hashes
    for i := range v.Array {
      validate_types(v.Array[i])
    }
    for i := range v.Hash {
      validate_types(v.Hash[i])
    }
    //////////////////////////////////////

    /////////////////////////////////////////////

  }

}
