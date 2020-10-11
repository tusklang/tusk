package interpreter

import (
	"fmt"
	"os"

	. "github.com/tusklang/tusk/lang/types"
	. "github.com/tusklang/tusk/native"
)

type overload_after struct {
	name string
	val  *TuskType
}

//export FillIns
func FillIns(instance *Instance, compiledVars map[string]*TuskType, dirname string, args []string) map[string]*TuskVar {
	globals := make(map[string]*TuskVar)
	var dirnameTuskStr TuskString
	dirnameTuskStr.FromGoType(dirname)
	var dirnameTuskType TuskType = dirnameTuskStr

	globals["__dirname"] = &TuskVar{
		Name:  "__dirname",
		Value: &dirnameTuskType,
	}

	var argv = make([]*TuskType, len(args))

	for k, v := range args {
		var kastr TuskString
		kastr.FromGoType(v)
		var tusktype TuskType = kastr
		argv[k] = &tusktype
	}

	var arr TuskType = TuskArray{
		Array:  argv,
		Length: uint64(len(args)),
	}

	globals["argv"] = &TuskVar{
		Name:  "argv",
		Value: &arr,
	}

	for k, v := range compiledVars {
		var global = TuskVar{
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

func RunInterpreter(compiledVars map[string]*TuskType, cli_params CliParams) {

	var instance Instance
	instance.Params = cli_params
	globals := FillIns(&instance, compiledVars, cli_params.Directory, os.Args)

	if _, exists := globals["main"]; !exists {
		fmt.Println("Given program has no entry point")
		os.Exit(1)
	} else {

		switch (*globals["main"].Value).(type) {
		case TuskFunc:
			main := globals["main"]

			if len((*main.Value).(TuskFunc).Overloads[0].Params) != 0 {
				fmt.Println("Invalid entry point")
				os.Exit(1)
			}

			calledP, e := Interpreter(&instance, (*main.Value).(TuskFunc).Overloads[0].Body, []string{"at the entry caller"}, 0, nil, false)
			if e != nil {
				e.Print()
				os.Exit(1)
			}
			WaitAllThreads() //finish up any remaining threads

			//get the exit value
			calledE := calledP.Exp
			if calledE == nil {
				return
			}
			called := *calledE
			////////////////////

			var exitType int64

			switch called.(type) {
			case TuskNumber:
				exitType = int64(called.(TuskNumber).ToGoType())
			}

			os.Exit(int(exitType))
		default:
			fmt.Println("Entry point was not given as a function")
			os.Exit(1)
		}
	}
}
