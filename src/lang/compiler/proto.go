package compiler

import . "lang/types"

var protos []string

//determine if the actions have not global prototypes
/*
  var main: fn() {
    var test: proto { ; would cause an error
    }
  }
*/
func has_non_global_prototypes(actions []Action, firstLayer bool) {

  for _, v := range actions {

    if v.Type == "proto" && !firstLayer {
      compilerErr("Prototypes can only be made at the global scope", v.File, v.Line)
    }

    if v.Type == "proto" {

      for i := range v.Static {
        has_non_global_prototypes(v.Static[i], false)
      }
      for i := range v.Instance {
        has_non_global_prototypes(v.Instance[i], false)
      }

      continue
    }
    if v.Type == "function" {
      has_non_global_prototypes(v.Value.(OmmFunc).Body, false)
      continue
    }

    if v.Type == "var" {
      has_non_global_prototypes(v.ExpAct, firstLayer)
      continue
    }

    //perform checker on all of the sub actions
    has_non_global_prototypes(v.ExpAct, false)
    has_non_global_prototypes(v.First, false)
    has_non_global_prototypes(v.Second, false)
    has_non_global_prototypes(v.Degree, false)

    //also do it for the arrays, hashes, and sub protos
    for _, i := range v.Array {
      has_non_global_prototypes(i, false)
    }
    for _, i := range v.Hash {
      has_non_global_prototypes(i, false)
    }
    for _, i := range v.Static {
      has_non_global_prototypes(i, false)
    }
    for _, i := range v.Instance {
      has_non_global_prototypes(i, false)
    }
    //////////////////////////////////////

    ///////////////////////////////////////////

  }
}

//put the proto names in the "types" slice
func put_proto_types(actions []Action) {

  for k, v := range actions {
    if v.Type == "var" && v.ExpAct[0].Type == "proto" {
      types = append(types, v.Name[1:])
      actions[k].ExpAct[0].Name = v.Name
    }
  }

}
