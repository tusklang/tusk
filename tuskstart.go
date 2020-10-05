package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tusklang/tools"
	"github.com/tusklang/tusk/lang/interpreter"
	. "github.com/tusklang/tusk/lang/types"

	"github.com/tusklang/tusk/lang/compiler"
)

var ver = flag.Bool("ver", false, "Get Tusk suite version")

func init() {
	flag.Usage = tools.Usagef("Tusk")
}

func main() {
	flag.Parse()

	if *ver {
		fmt.Printf("Tusk Beta %d.%d.%d", tools.TuskMajor, tools.TuskMinor, tools.TuskBug)
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

	acts, e := compiler.Compile(cli_params)

	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}

	os.Args = os.Args[1:] //remove the `tusk` <file.tusk>

	interpreter.RunInterpreter(acts, cli_params)
}
