package oat

import (
	"fmt"
	"os"

	"github.com/omm-lang/omm/lang/interpreter"
	"github.com/omm-lang/omm/lang/types"
	oatenc "github.com/omm-lang/omm/oat/encoding"
)

func Run(params types.CliParams) {
	d, e := oatenc.OatDecode(params.Name, 0)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	interpreter.RunInterpreter(d, params)
}
