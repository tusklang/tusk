package compiler

import (
	"fmt"
	"os"

	"github.com/tusklang/tusk/lang/types"

	"github.com/tusklang/tusk/lang/interpreter"
)

func Run(params types.CliParams) {
	variables, ce := Compile(params)

	if ce != nil {
		fmt.Println(ce.Error())
		os.Exit(1)
	}

	os.Args = os.Args[1:] //remove the `tusk` <file.tusk>

	interpreter.RunInterpreter(variables, params)
}
