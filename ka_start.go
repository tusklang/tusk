package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	. "ka/lang/types"
	"kore"

	"ka/lang/compiler"
)

var ver = flag.Bool("ver", false, "Get Ka suite version")

func init() {
	flag.Usage = kore.Usagef("Ka")
}

func main() {
	flag.Parse()

	if *ver {
		fmt.Printf("Kore Beta %d.%d.%d", suite.KaSuiteMajor, suite.KaSuiteMinor, suite.KaSuiteBug)
		os.Exit(0)
	}

	var cli_params CliParams

	if len(flag.Arg(0)) != 0 && flag.Arg(0)[0] == '-' {
		fmt.Println("Error, must list the filename as the first parameter")
		os.Exit(1)
	}

	cli_params.Name = flag.Arg(0)

	dirname, _ := os.Executable()

	cli_params.KaDirname = filepath.Dir(dirname)

	//set the working directory
	os.Chdir(cli_params.Directory)

	compiler.Run(cli_params)
}
