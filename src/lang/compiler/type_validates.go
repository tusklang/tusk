package compiler

import . "lang/types"

var types = []string{ "string", "rune", "number", "bool", "hash", "array" }

func validate_types(actions []Action) CompileErr {

  var e CompileErr

  //function to make sure the typecasts

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
      e = validate_types(v.Value.(OmmFunc).Body)
      if e != nil {
        return e
      }
      continue
    }
    if v.Type == "proto" {

      for i := range v.Value.(OmmProto).Static {
        var val = *v.Value.(OmmProto).Static[i]

        if val.Type() == "function" {
          e = validate_types(val.(OmmFunc).Body)
          if e != nil {
            return e
          }
        }
      }
      for i := range v.Value.(OmmProto).Instance {
        var val = *v.Value.(OmmProto).Instance[i]

        if val.Type() == "function" {
          e = validate_types(val.(OmmFunc).Body,)
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
    /////////////////////////////////////////////

  }

  return nil
}
