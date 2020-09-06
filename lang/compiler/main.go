package compiler

import (
	"fmt"
	"os"

	"github.com/omm-lang/omm/lang/types"

	"github.com/omm-lang/omm/lang/interpreter"
)

func Run(params types.CliParams) {

	fileName := params.Name

	included = append(included, fileName)

	variables, ce := Compile(params)

	if ce != nil {
		fmt.Println(ce.Error())
		os.Exit(1)
	}

	os.Args = os.Args[1:] //remove the `omm` <file.omm>

	interpreter.RunInterpreter(variables, params)
}
