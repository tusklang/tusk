package compiler

import "ka/lang/types"

func getvars(actions []types.Action) (map[string]*types.KaType, error) {

	var vars = make(map[string]*types.KaType)

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
				var f types.KaType = types.KaFunc{}
				vars[v.Name] = &f
			}

			if (*vars[v.Name]).Type() != "function" {
				return nil, makeCompilerErr(v.Name[1:]+" is not a function", v.File, v.Line)
			}

			tmp := (*vars[v.Name]).(types.KaFunc)
			tmp.Overloads = append(tmp.Overloads, v.ExpAct[0].Value.(types.KaFunc).Overloads...)
			var katype types.KaType = tmp
			vars[v.Name] = &katype
			continue
		}

		vars[v.Name] = &v.ExpAct[0].Value
	}

	return vars, nil
}
