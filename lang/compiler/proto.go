package compiler

import (
	. "tusk/lang/types"
)

var protos []string

//determine if the actions have not global prototypes
/*
  var main: fn() {
    var test: proto { ; would cause an error

    }
  }
*/
func has_non_global_prototypes(actions []Action, firstLayer bool) error {

	var e error

	for _, v := range actions {

		if v.Type == "proto" && !firstLayer {
			return makeCompilerErr("Prototypes can only be made at the global scope", v.File, v.Line)
		}

		if v.Type == "proto" {

			for i := range v.Value.(TuskProto).Static {
				var val = *v.Value.(TuskProto).Static[i]

				if val.Type() == "function" {
					e = has_non_global_prototypes([]Action{Action{
						Type:  "function",
						Value: val,
					}}, false)
					if e != nil {
						return e
					}
				}
			}
			for i := range v.Value.(TuskProto).Instance {
				var val = *v.Value.(TuskProto).Instance[i]

				if val.Type() == "function" {
					e = has_non_global_prototypes([]Action{Action{
						Type:  "function",
						Value: val,
					}}, false)
					if e != nil {
						return e
					}
				}
			}

			continue
		}
		if v.Type == "function" {
			e = has_non_global_prototypes(v.Value.(TuskFunc).Overloads[0].Body, false)
			if e != nil {
				return e
			}
			continue
		}

		if v.Type == "var" {
			e = has_non_global_prototypes(v.ExpAct, firstLayer)
			if e != nil {
				return e
			}
			continue
		}

		//perform checker on all of the sub actions
		e = has_non_global_prototypes(v.ExpAct, false)
		if e != nil {
			return e
		}
		e = has_non_global_prototypes(v.First, false)
		if e != nil {
			return e
		}
		e = has_non_global_prototypes(v.Second, false)
		if e != nil {
			return e
		}

		//also do it for the (runtime) arrays and hashes
		for i := range v.Array {
			e = has_non_global_prototypes(v.Array[i], false)
			if e != nil {
				return e
			}
		}
		for i := range v.Hash {
			e = has_non_global_prototypes(v.Hash[i][0], false)
			e = has_non_global_prototypes(v.Hash[i][1], false)
			if e != nil {
				return e
			}
		}
		////////////////////////////////////////////////

		///////////////////////////////////////////

	}

	return e
}

//put the proto names in the "types" slice
func put_proto_types(actions []Action) {

	for k, v := range actions {
		if v.Type == "var" && len(v.ExpAct) != 0 && v.ExpAct[0].Type == "proto" {
			validtypes = append(validtypes, v.Name[1:])

			var proto = v.ExpAct[0].Value.(TuskProto)
			proto.ProtoName = v.Name
			actions[k].ExpAct[0].Value = proto
		}
	}

}
