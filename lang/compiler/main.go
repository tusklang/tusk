package compiler

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/omm-lang/omm/lang/types"

	. "github.com/omm-lang/omm/lang/interpreter"
)

var included = []string{} //list of the imported files from omm

func Run(params CliParams) {

	fileName := params.Name

	included = append(included, fileName)

	file, e := ioutil.ReadFile(fileName)

	if e != nil {
		fmt.Println("Could not find", fileName)
		os.Exit(1)
	}

	variables, ce := Compile(string(file), fileName, params)

	if ce != nil {
		ce.Print()
	}

	RunInterpreter(variables, params)
}
