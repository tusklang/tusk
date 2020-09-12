package compiler

import (
	"fmt"
	"os"

	"omm/lang/types"

	"omm/lang/interpreter"
)

func Run(params types.CliParams) {
	variables, ce := Compile(params)

	if ce != nil {
		fmt.Println(ce.Error())
		os.Exit(1)
	}

	os.Args = os.Args[1:] //remove the `omm` <file.omm>

	interpreter.RunInterpreter(variables, params)
}
