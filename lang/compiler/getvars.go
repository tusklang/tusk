package compiler

import "github.com/tusklang/tusk/lang/types"

func mergeproto(m1, m2 map[string]*types.TuskType) {
	for k, v := range m2 {
		m1[k] = v
	}
}

func mergeprotoaccess(m1, m2 map[string][]string) {
	for k, v := range m2 {
		m1[k] = v
	}
}

func getvars(actions []types.Action) (map[string]*types.TuskType, error) {

	var vars = make(map[string]*types.TuskType)

	for _, v := range actions {
		if v.Type != "var" && v.Type != "declare" && v.Type != "ovld" { //if it is not an assigner or overloader, it must be an error
			return nil, makeCompilerErr("Cannot have anything but a variable declaration or overloader outside of a function", v.File, v.Line)
		}

		if v.Type == "declare" {
			var undef types.TuskType = types.TuskUndef{}
			vars[v.Name] = &undef
			continue
		}

		if v.ExpAct[0].Value == nil { //if it has no value (meaning a compound type)
			return nil, makeCompilerErr("Cannot have compound types at the global scope", v.File, v.Line)
		}

		if v.Type == "ovld" {
			if _, exists := vars[v.Name]; !exists || (*vars[v.Name]).Type() == "none" { //if it does not exist yet, declare undefined yet
				var f types.TuskType = types.TuskFunc{}
				vars[v.Name] = &f
			}

			if (*vars[v.Name]).Type() != "function" {
				return nil, makeCompilerErr(v.Name+" is not a function", v.File, v.Line)
			}

			tmp := (*vars[v.Name]).(types.TuskFunc)
			tmp.Overloads = append(tmp.Overloads, v.ExpAct[0].Value.(types.TuskFunc).Overloads...)
			var tusktype types.TuskType = tmp
			vars[v.Name] = &tusktype
			continue
		}

		if v.ExpAct[0].Value.Type() == "prototype" {
			if proto, exists := vars[v.Name]; exists {

				/* If two protos with the same name are declares, merge them (for namespacing purposes)
				var example = proto {
					var a
				}

				var example = proto {
					var b
				}

				becomes

				var example = proto {
					var a
					var b
				}
				*/

				mergeproto((*proto).(types.TuskProto).Static, v.ExpAct[0].Value.(types.TuskProto).Static)
				mergeproto((*proto).(types.TuskProto).Instance, v.ExpAct[0].Value.(types.TuskProto).Instance)
				mergeprotoaccess((*proto).(types.TuskProto).AccessList, v.ExpAct[0].Value.(types.TuskProto).AccessList)
				continue
			}
		}

		vars[v.Name] = &v.ExpAct[0].Value
	}

	return vars, nil
}
