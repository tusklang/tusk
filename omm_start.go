package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	suite "omm-suite"
	. "omm/lang/types"

	"omm/lang/compiler"
)

var ver = flag.Bool("ver", false, "Get Omm suite version")

func init() {
	flag.Usage = suite.Usagef("Omm")
}

func main() {
	flag.Parse()

	if *ver {
		fmt.Printf("Omm Beta %d.%d.%d", suite.OmmSuiteMajor, suite.OmmSuiteMinor, suite.OmmSuiteBug)
		os.Exit(0)
	}

	var cli_params CliParams

	if len(flag.Arg(0)) != 0 && flag.Arg(0)[0] == '-' {
		fmt.Println("Error, must list the filename as the first parameter")
		os.Exit(1)
	}

	cli_params.Name = flag.Arg(0)

	dirname, _ := os.Executable()

	cli_params.OmmDirname = filepath.Dir(dirname)

	//set the working directory
	os.Chdir(cli_params.Directory)

	compiler.Run(cli_params)
}
