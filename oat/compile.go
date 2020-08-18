package oat

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// import "encoding/gob"
	// import "bytes"

	. "github.com/omm-lang/omm/lang/types"

	"github.com/omm-lang/omm/lang/compiler"

	. "github.com/omm-lang/omm/oat/encoding"
)

//compiler

//export Compile
func Compile(params CliParams) {
	fileName := params.Name

	file, e := ioutil.ReadFile(fileName)

	if e != nil {
		fmt.Println("Could not find file:", fileName)
		os.Exit(1)
	}

	var compileall = false
	if strings.HasSuffix(fileName, "*") || strings.HasSuffix(fileName, "*/") {
		compileall = true
		fileName = "main.omm"
	}

	compiler.Ommbasedir = params.OmmDirname
	vars, ce := compiler.Compile(string(file), fileName, compileall, true)
	compiler.Ommbasedir = "" //reset Ommbasedir

	if ce != nil {
		ce.Print()
	}

	OatEncode(params.Output, vars)
}
