package compiler

import (
	"path/filepath"
	"strings"

	"github.com/tusklang/tusk/lang/types"
)

//this file allows for globals to have specified access

func globalAccess(actions []types.Action) ([]types.Action, map[string][]string) {
	var nActions []types.Action //list of new actions (without access)
	var accessl = make(map[string][]string)

	for k, v := range actions {

		if v.Type == "access" {
			//an access list
			if k+1 < len(actions) { //if the next action exists
				v.Value.Range(func(kk, vv *types.TuskType) (types.Returner, *types.TuskError) {
					//get the file
					name := (*vv).Format()

					//thisf is a macro for the current file
					if name == "thisf" {
						name = v.File
					} else if strings.HasPrefix(name, "curdir") {
						name = filepath.Dir(v.File) + strings.TrimPrefix(name, "curdir")
					}

					accessl[actions[k+1].Name] = append(accessl[actions[k+1].Name], name)
					return types.Returner{}, nil
				})
			} //otherwise, just ignore it
			continue
		}

		nActions = append(nActions, v)
	}

	return nActions, accessl
}
