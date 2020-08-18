package compiler

import . "github.com/omm-lang/omm/lang/types"

func getFunctionVarRefs(actions []Action) {

	//this is to support closures
	//this gets all of the variables that this function references, prevent them from being garbage collected

	for k, v := range actions {

		if v.Type == "function" {
			var fn = v.Value.(OmmFunc).Overloads[0]
			fn.VarRefs = putFunctionVarRefs(fn.Body)
			getFunctionVarRefs(fn.Body)
			actions[k].Value.(OmmFunc).Overloads[0] = fn
			continue
		}

		//perform checkvars on all of the sub actions
		getFunctionVarRefs(v.ExpAct)
		getFunctionVarRefs(v.First)
		getFunctionVarRefs(v.Second)

		//also do it for the (runtime) arrays and hashes
		for i := range v.Array {
			getFunctionVarRefs(actions[k].Array[i])
		}
		for i := range v.Hash {
			getFunctionVarRefs(actions[k].Hash[i][0])
			getFunctionVarRefs(actions[k].Hash[i][1])
		}
		////////////////////////////////////////////////

		/////////////////////////////////////////////

	}

}

func putFunctionVarRefs(body []Action) []string {

	var names []string

	for _, v := range body {

		if v.Type == "function" {

			var fn = v.Value.(OmmFunc).Overloads[0]

			fn.VarRefs = putFunctionVarRefs(fn.Body)

			fn.VarRefs = append(fn.VarRefs)

			putFunctionVarRefs(fn.Body)
			continue
		}

		if v.Type == "variable" {
			names = append(names, v.Name)
		}

		//perform checkvars on all of the sub actions
		putFunctionVarRefs(v.ExpAct)
		putFunctionVarRefs(v.First)
		putFunctionVarRefs(v.Second)

		//also do it for the (runtime) arrays and hashes
		for i := range v.Array {
			putFunctionVarRefs(v.Array[i])
		}
		for i := range v.Hash {
			putFunctionVarRefs(v.Hash[i][0])
			putFunctionVarRefs(v.Hash[i][1])
		}
		////////////////////////////////////////////////

		/////////////////////////////////////////////

	}

	return names
}
