package compiler

import (
	"fmt"
	"os"

	"ka/lang/types"

	"ka/lang/interpreter"
)

func Run(params types.CliParams) {
	variables, ce := Compile(params)

	if ce != nil {
		fmt.Println(ce.Error())
		os.Exit(1)
	}

	os.Args = os.Args[1:] //remove the `ka` <file.ka>

	interpreter.RunInterpreter(variables, params)
}
