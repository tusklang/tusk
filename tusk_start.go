package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	. "tusk/lang/types"
	"kore"

	"tusk/lang/compiler"
)

var ver = flag.Bool("ver", false, "Get Tusk suite version")

func init() {
	flag.Usage = kore.Usagef("Tusk")
}

func main() {
	flag.Parse()

	if *ver {
		fmt.Printf("Kore Beta %d.%d.%d", kore.KoreMajor, kore.KoreMinor, kore.KoreBug)
		os.Exit(0)
	}

	var cli_params CliParams

	if len(flag.Arg(0)) != 0 && flag.Arg(0)[0] == '-' {
		fmt.Println("Error, must list the filename as the first parameter")
		os.Exit(1)
	}

	cli_params.Name = flag.Arg(0)

	dirname, _ := os.Executable()

	cli_params.TuskDirname = filepath.Dir(dirname)

	//set the working directory
	os.Chdir(cli_params.Directory)

	compiler.Run(cli_params)
}
