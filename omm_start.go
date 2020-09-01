package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/omm-lang/omm/lang/types"

	suite "github.com/omm-lang/omm-suite"
	"github.com/omm-lang/omm/lang/compiler"
)

var cwd = flag.String("cwd", "", "set cwd")

func defaults(cli_params *CliParams) {
	(*cli_params).Prec = 30
	(*cli_params).Name = ""
	(*cli_params).Directory = ""
}

func main() {
	flag.Parse()

	var cli_params CliParams

	if len(os.Args) <= 2 {
		fmt.Println("Error, no input file was given")
		os.Exit(1)
	}

	defaults(&cli_params)

	cli_params.Directory = *cwd
	cli_params.Name = flag.Arg(0)

	dirname, _ := os.Executable()

	cli_params.OmmDirname = filepath.Dir(dirname)

	//set the working directory
	os.Chdir(cli_params.Directory)

	for i := 1; i < len(flag.Args()); i++ {

		v := flag.Arg(i)

		if strings.HasPrefix(v, "--") {

			switch v[2:] {

			case "version":
				fmt.Printf("Omm Beta %d.%d.%d", suite.OmmSuiteMajor, suite.OmmSuiteMinor, suite.OmmSuiteBug)
				os.Exit(0)
			}

		} else if strings.HasPrefix(v, "-") {

			switch v[1:] {

			case "v":
				fmt.Printf("Omm Beta %d.%d.%d", suite.OmmSuiteMajor, suite.OmmSuiteMinor, suite.OmmSuiteBug)
				os.Exit(0)
			case "prec":
				tempprec, _ := strconv.ParseUint(flag.Arg(i+1), 10, 64)
				cli_params.Prec = tempprec
				i += 2
			default:
				fmt.Println("Warning, there is no cli parameter named", v)
				i++
			}
		}
	}

	compiler.Run(cli_params)
}
