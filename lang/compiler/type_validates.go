package compiler

import (
	"tusk/lang/interpreter"
	. "tusk/lang/types"
)

var validtypes = []string{"string", "rune", "number", "bool", "hash", "array", "function", "none", "any"}

func getnativetypes() {
	for _, v := range interpreter.Native { //append all of the types of the native
		validtypes = append(validtypes, (*v).Type())
	}
}

func validate_types(actions []Action) error {

	var e error

	//function to make sure the typecasts and paramlists have valid types

	for _, v := range actions {

		if v.Type == "cast" {
			for _, t := range validtypes {
				if t == v.Name { //if the type exists, do not throw an error
					goto cast_noErr
				}
			}
			return makeCompilerErr("\""+v.Name+"\" is not a type", v.File, v.Line)
		cast_noErr:
		}

		if v.Type == "function" {

			//check the parameter list
			for _, vv := range v.Value.(TuskFunc).Overloads[0].Types {
				for _, t := range validtypes {
					if t == vv {
						goto plist_noErr
					}
				}

				return makeCompilerErr("\""+vv+"\" is not a type", v.File, v.Line)

			plist_noErr:
			}

			e = validate_types(v.Value.(TuskFunc).Overloads[0].Body)
			if e != nil {
				return e
			}
			continue
		}
		if v.Type == "proto" {

			for i := range v.Value.(TuskProto).Static {
				var val = *v.Value.(TuskProto).Static[i]

				if val.Type() == "function" {
					e = validate_types([]Action{Action{
						Type:  "function",
						Value: val,
					}})
					if e != nil {
						return e
					}
				}
			}
			for i := range v.Value.(TuskProto).Instance {
				var val = *v.Value.(TuskProto).Instance[i]

				if val.Type() == "function" {
					e = validate_types([]Action{Action{
						Type:  "function",
						Value: val,
					}})
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
			e = validate_types(v.Hash[i][0])
			e = validate_types(v.Hash[i][1])
			if e != nil {
				return e
			}
		}
		////////////////////////////////////////////////

		/////////////////////////////////////////////

	}

	return nil
}
