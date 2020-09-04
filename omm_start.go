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

	var filenamei int
	for flag.Arg(filenamei) != "" && flag.Arg(filenamei)[0] == '-' {
		filenamei++ //only inside the block for formatting
	}

	if flag.Arg(filenamei) == "" {
		fmt.Println("Error, no input file was given")
		os.Exit(1)
	}

	cli_params.Name = flag.Arg(0)

	dirname, _ := os.Executable()

	cli_params.OmmDirname = filepath.Dir(dirname)

	//set the working directory
	os.Chdir(cli_params.Directory)

	compiler.Run(cli_params)
}
