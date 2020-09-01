package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	suite "github.com/omm-lang/omm-suite"
	. "github.com/omm-lang/omm/lang/types"

	"github.com/omm-lang/omm/lang/compiler"
)

var cwd = flag.String("cwd", "", "set the current working directory (automatically placed by the shell/pwsh script)")
var ver = flag.String("ver", "", "get Omm suite version")
var prec = flag.Uint64("prec", 20, "set the precision of an Omm instance")

func init() {
	flag.Usage = suite.Usagef("Omm")
}

func main() {
	flag.Parse()

	var cli_params CliParams

	if len(os.Args) <= 2 {
		fmt.Println("Error, no input file was given")
		os.Exit(1)
	}

	cli_params.Directory = *cwd
	cli_params.Name = flag.Arg(0)
	cli_params.Prec = *prec

	dirname, _ := os.Executable()

	cli_params.OmmDirname = filepath.Dir(dirname)

	//set the working directory
	os.Chdir(cli_params.Directory)

	compiler.Run(cli_params)
}
