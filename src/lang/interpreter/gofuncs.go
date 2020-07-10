package interpreter

//all of the gofuncs
//functions written in go that are used by omm

import . "lang/types"

//export GoFuncs
var GoFuncs = map[string]func(args [][]Action, cli_params CliParams, vars map[string]Variable, this_vals []Action, dir string) OmmType {}
