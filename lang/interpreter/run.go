package interpreter

import (
	"fmt"
	"os"

	. "ka/lang/types"
)

type overload_after struct {
	name string
	val  *KaType
}

//export FillIns
func FillIns(instance *Instance, compiledVars map[string]*KaType, dirname string, args []string) map[string]*KaVar {
	globals := make(map[string]*KaVar)
	var dirnameKaStr KaString
	dirnameKaStr.FromGoType(dirname)
	var dirnameKaType KaType = dirnameKaStr

	globals["$__dirname"] = &KaVar{
		Name:  "$__dirname",
		Value: &dirnameKaType,
	}

	var argv = make([]*KaType, len(args))

	for k, v := range args {
		var kastr KaString
		kastr.FromGoType(v)
		var katype KaType = kastr
		argv[k] = &katype
	}

	var arr KaType = KaArray{
		Array:  argv,
		Length: uint64(len(args)),
	}

	globals["$argv"] = &KaVar{
		Name:  "$argv",
		Value: &arr,
	}

	for k, v := range compiledVars {
		var global = KaVar{
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

func RunInterpreter(compiledVars map[string]*KaType, cli_params CliParams) {

	var instance Instance
	instance.Params = cli_params
	globals := FillIns(&instance, compiledVars, cli_params.Directory, os.Args)

	if _, exists := globals["$main"]; !exists {
		fmt.Println("Given program has no entry point/main function")
		os.Exit(1)
	} else {

		switch (*globals["$main"].Value).(type) {
		case KaFunc:
			main := globals["$main"]

			calledP := Interpreter(&instance, (*main.Value).(KaFunc).Overloads[0].Body, []string{"at the entry caller"}, 0, nil, false).Exp
			WaitAllThreads() //finish up any remaining threads

			if calledP == nil {
				os.Exit(0)
			}

			called := *calledP //dereference the called ptr (to get the value)

			var exitType int64

			switch called.(type) {
			case KaNumber:
				exitType = int64(called.(KaNumber).ToGoType())
			}

			os.Exit(int(exitType))
		default:
			fmt.Println("Entry point was not given as a function")
			os.Exit(1)
		}
	}
}
