package compiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/omm-lang/omm/lang/types"

	. "github.com/omm-lang/omm/lang/interpreter"
)

var included = []string{} //list of the imported files from omm

//export Ommbasedir
var Ommbasedir string //directory of the omm installation

//export Run
func Run(params CliParams) {

	fileName := params.Name

	var compileall = false
	if strings.HasSuffix(fileName, "*") || strings.HasSuffix(fileName, "*/") {
		compileall = true
		fileName = "main.omm"
	}

	included = append(included, fileName)

	file, e := ioutil.ReadFile(fileName)

	if e != nil {
		fmt.Println("Could not find", fileName)
		os.Exit(1)
	}

	Ommbasedir = params.OmmDirname
	variables, ce := Compile(string(file), fileName, compileall, true)
	Ommbasedir = "" //reset Ommbasedir

	if ce != nil {
		ce.Print()
	}

	RunInterpreter(variables, params)
}
