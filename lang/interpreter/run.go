package interpreter

import (
	"fmt"
	"os"

	. "github.com/omm-lang/omm/lang/types"
)

type overload_after struct {
	name string
	val  *OmmType
}

//export FillIns
func FillIns(instance *Instance, compiledVars map[string]*OmmType, dirname string, args []string) map[string]*OmmVar {
	globals := make(map[string]*OmmVar)
	var dirnameOmmStr OmmString
	dirnameOmmStr.FromGoType(dirname)
	var dirnameOmmType OmmType = dirnameOmmStr

	globals["$__dirname"] = &OmmVar{
		Name:  "$__dirname",
		Value: &dirnameOmmType,
	}

	var argv = make([]*OmmType, len(args))

	for k, v := range args {
		var ommstr OmmString
		ommstr.FromGoType(v)
		var ommtype OmmType = ommstr
		argv[k] = &ommtype
	}

	var arr OmmType = OmmArray{
		Array:  argv,
		Length: uint64(len(args)),
	}

	globals["$argv"] = &OmmVar{
		Name:  "$argv",
		Value: &arr,
	}

	for k, v := range compiledVars {
		var global = OmmVar{
			Name:  k,
			Value: v,
		}
		globals[k] = &global
	}

	//also allocate to the locals
	for k, v := range globals {
		instance.Allocate(k, v.Value)
	}

	for k, v := range Native { //allocate all of the native
		instance.Allocate(k, v)
	}

	return globals
}

func RunInterpreter(compiledVars map[string]*OmmType, cli_params CliParams) {

	var instance Instance
	instance.Params = cli_params
	globals := FillIns(&instance, compiledVars, cli_params.Directory, os.Args)

	if _, exists := globals["$main"]; !exists {
		fmt.Println("Given program has no entry point/main function")
		os.Exit(1)
	} else {

		switch (*globals["$main"].Value).(type) {
		case OmmFunc:
			main := globals["$main"]

			calledP := Interpreter(&instance, (*main.Value).(OmmFunc).Overloads[0].Body, []string{"at the entry caller"}, 0, nil, false).Exp
			WaitAllThreads() //finish up any remaining threads

			if calledP == nil {
				os.Exit(0)
			}

			called := *calledP //dereference the called ptr (to get the value)

			var exitType int64

			switch called.(type) {
			case OmmNumber:
				exitType = int64(called.(OmmNumber).ToGoType())
			}

			os.Exit(int(exitType))
		default:
			fmt.Println("Entry point was not given as a function")
			os.Exit(1)
		}
	}
}
