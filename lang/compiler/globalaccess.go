package compiler

import (
	"github.com/tusklang/tusk/lang/types"
)

//this file allows for globals to have specified access

type accessStruct struct {
	path   string
	access []string
}

func globalAccess(actions []types.Action) ([]types.Action, map[string]accessStruct) {
	var nActions []types.Action //list of new actions (without access)
	var accessl = make(map[string]accessStruct)

	for k, v := range actions {

		if v.Type == "access" {
			//an access list
			if k+1 < len(actions) { //if the next action exists

				var constructAccess = accessStruct{
					path:   v.File,
					access: []string{},
				}

				v.Value.Range(func(kk, vv *types.TuskType) (types.Returner, *types.TuskError) {
					//get the file
					name := (*vv).Format()

					if name == "thisf" {
						name = v.File
					}

					constructAccess.access = append(constructAccess.access, name)
					return types.Returner{}, nil
				})
				accessl[actions[k+1].Name] = constructAccess
			} //otherwise, just ignore it
			continue
		}

		nActions = append(nActions, v)
	}

	return nActions, accessl
}
