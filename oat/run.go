package oat

import (
	"fmt"
	"os"

	"github.com/omm-lang/omm/lang/interpreter"

	. "github.com/omm-lang/omm/oat/encoding"

	. "github.com/omm-lang/omm/lang/types"
)

//export Run
func Run(params CliParams) {
	d, e := OatDecode(params.Name, 0)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	interpreter.RunInterpreter(d, params)
}
