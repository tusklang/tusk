package compiler

import . "lang/types"

var types = []string{ "string", "rune", "number", "bool", "hash", "array", "function", "any" }

func validate_types(actions []Action) CompileErr {

  var e CompileErr

  //function to make sure the typecasts and paramlists have valid types

  for _, v := range actions {

    if v.Type == "cast" {
      for _, t := range types {
        if t == v.Name { //if the type exists, do not throw an error
          goto cast_noErr
        }
      }
      return makeCompilerErr("\"" + v.Name + "\" is not a type", v.File, v.Line)
      cast_noErr:
    }

    if v.Type == "function" {

      //check the parameter list
      for _, vv := range v.Value.(OmmFunc).Overloads[0].Types {
        for _, t := range types {
          if t == vv {
            goto plist_noErr
          }
        }

        return makeCompilerErr("\"" + vv + "\" is not a type", v.File, v.Line)

        plist_noErr:
      }

      e = validate_types(v.Value.(OmmFunc).Overloads[0].Body)
      if e != nil {
        return e
      }
      continue
    }
    if v.Type == "proto" {

      for i := range v.Static {
        var val = v.Static[i][0]

        if val.Type == "function" {
          e = validate_types([]Action{ val })
          if e != nil {
            return e
          }
        }
      }
      for i := range v.Instance {
        var val = v.Instance[i][0]

        if val.Type == "function" {
          e = validate_types([]Action{ val })
          if e != nil {
            return e
          }
        }
      }

      continue
    }

    //perform checkvars on all of the sub actions
    e = validate_types(v.ExpAct)
    if e != nil {
      return e
    }
    e = validate_types(v.First)
    if e != nil {
      return e
    }
    e = validate_types(v.Second)
    if e != nil {
      return e
    }

    //also do it for the (runtime) arrays and hashes
    for i := range v.Array {
      e = validate_types(v.Array[i])
      if e != nil {
        return e
      }
    }
    for i := range v.Hash {
      e = validate_types(v.Hash[i])
      if e != nil {
        return e
      }
    }
    ////////////////////////////////////////////////

    /////////////////////////////////////////////

  }

  return nil
}
