package compiler

import "tusk/lang/types"

func getvars(actions []types.Action) (map[string]*types.TuskType, error) {

	var vars = make(map[string]*types.TuskType)

	for _, v := range actions {
		if v.Type != "var" && v.Type != "declare" && v.Type != "ovld" { //if it is not an assigner or overloader, it must be an error
			return nil, makeCompilerErr("Cannot have anything but a variable declaration or overloader outside of a function", v.File, v.Line)
		}

		if v.Type == "declare" {
			continue
		}

		if v.ExpAct[0].Value == nil { //if it has no value (meaning a compound type)
			return nil, makeCompilerErr("Cannot have compound types at the global scope", v.File, v.Line)
		}

		if v.Type == "ovld" {
			if _, exists := vars[v.Name]; !exists { //if it does not exist yet, declare undefined yet
				var f types.TuskType = types.TuskFunc{}
				vars[v.Name] = &f
			}

			if (*vars[v.Name]).Type() != "function" {
				return nil, makeCompilerErr(v.Name[1:]+" is not a function", v.File, v.Line)
			}

			tmp := (*vars[v.Name]).(types.TuskFunc)
			tmp.Overloads = append(tmp.Overloads, v.ExpAct[0].Value.(types.TuskFunc).Overloads...)
			var tusktype types.TuskType = tmp
			vars[v.Name] = &tusktype
			continue
		}

		vars[v.Name] = &v.ExpAct[0].Value
	}

	return vars, nil
}
