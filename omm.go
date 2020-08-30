package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/omm-lang/omm/lang/types"

	"github.com/omm-lang/framework"
	"github.com/omm-lang/omm/lang/compiler"
)

func defaults(cli_params *CliParams) {
	(*cli_params).Prec = 30
	(*cli_params).Name = ""
	(*cli_params).Directory = ""
}

func main() {

	args := os.Args

	var cli_params CliParams

	if len(args) <= 1 {
		fmt.Println("Error, no input file was given")
		os.Exit(1)
	}

	defaults(&cli_params)

	cli_params.Directory, _ = os.Getwd()
	cli_params.Name = args[1]

	dirname, _ := os.Executable()

	cli_params.OmmDirname = filepath.Dir(dirname)

	//set the working directory
	os.Chdir(cli_params.Directory)

	for i := 2; i < len(args); i++ {

		v := args[i]

		if strings.HasPrefix(v, "--") {

			switch v[2:] {

			case "version":
				fmt.Printf("Omm Beta %d.%d.%d", framework.OmmFrameworkMajor, framework.OmmFrameworkMinor, framework.OmmFrameworkBug)
				os.Exit(0)
			}

		} else if strings.HasPrefix(v, "-") {

			switch v[1:] {

			case "v":
				fmt.Printf("Omm Beta %d.%d.%d", framework.OmmFrameworkMajor, framework.OmmFrameworkMinor, framework.OmmFrameworkBug)
				os.Exit(0)
			case "prec":
				temp_prec, _ := strconv.ParseUint(args[i+1], 10, 64)
				cli_params.Prec = temp_prec
				i += 2
			default:
				fmt.Println("Warning, there is no cli parameter named", v)
				i++
			}
		}
	}

	compiler.Run(cli_params)
}
