package compiler

import (
	. "github.com/tusklang/tusk/lang/types"
)

func putFunctionVarRefs(body []Action) []string {

	var names []string

	for k, v := range body {

		if v.Type == "function" {

			var fn = v.Value.(TuskFunc)

			for k := range fn.Overloads {
				fn.Overloads[k].VarRefs = putFunctionVarRefs(fn.Overloads[k].Body)
				putFunctionVarRefs(fn.Overloads[k].Body)
			}

			body[k].Value = fn
			continue
		}

		if v.Type == "proto" {

			var proto = v.Value.(TuskProto)

			for k, v := range proto.Static {
				var passarr = []Action{Action{
					Type:  (*v).Type(),
					Value: *v,
				}}
				putFunctionVarRefs(passarr)
				proto.Static[k] = &passarr[0].Value
			}

			for k, v := range proto.Instance {
				var passarr = []Action{Action{
					Type:  (*v).Type(),
					Value: *v,
				}}
				putFunctionVarRefs(passarr)
				proto.Instance[k] = &passarr[0].Value
			}

			body[k].Value = proto
		}

		if v.Type == "variable" {
			names = append(names, v.Name)
		}

		//perform checkvars on all of the sub actions
		names = append(names, putFunctionVarRefs(v.ExpAct)...)
		names = append(names, putFunctionVarRefs(v.First)...)
		names = append(names, putFunctionVarRefs(v.Second)...)

		//also do it for the (runtime) arrays and hashes
		for i := range v.Array {
			names = append(names, putFunctionVarRefs(v.Array[i])...)
		}
		for i := range v.Hash {
			names = append(names, putFunctionVarRefs(v.Hash[i][0])...)
			names = append(names, putFunctionVarRefs(v.Hash[i][1])...)
		}
		////////////////////////////////////////////////

		/////////////////////////////////////////////

	}

	return names
}
