package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/omm-lang/omm/oat"

	. "github.com/omm-lang/omm/lang/types"

	"github.com/omm-lang/omm/lang/compiler"
)

func defaults(cli_params *CliParams, name string) {

	(*cli_params).Prec = 30

	if strings.HasSuffix(name, "*") || strings.HasSuffix(name, "*/") { //detect a directory compile
		(*cli_params).Output = "all.oat"
	} else if strings.LastIndex(name, ".") == -1 {
		(*cli_params).Output = name + ".oat"
	} else {
		(*cli_params).Output = name[:strings.LastIndex(name, ".")] + ".oat"
	}

	(*cli_params).Addon = "lang"
	(*cli_params).Name = ""
	(*cli_params).Directory = ""
}

func main() {

	args := os.Args

	var cli_params CliParams

	if len(args) == 0 {
		fmt.Println("Error, no input file was given")
		os.Exit(1)
	}

	defaults(&cli_params, args[1])

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
				fmt.Printf("Omm Beta %d.%d.%d", OMM_MAJOR, OMM_MINOR, OMM_BUG)
				os.Exit(0)
			default:
				cli_params.Addon = v[2:]
			}

		} else if strings.HasPrefix(v, "-") {

			switch v[1:] {

			case "v":
				fmt.Printf("Omm Beta %d.%d.%d", OMM_MAJOR, OMM_MINOR, OMM_BUG)
				os.Exit(0)
			case "c":
				cli_params.Addon = "compile"
			case "r":
				cli_params.Addon = "run"
			case "prec":
				temp_prec, _ := strconv.ParseUint(args[i+1], 10, 64)
				cli_params.Prec = temp_prec
				i += 2
			case "o":
				cli_params.Output = args[i+1]
				i++
			default:
				fmt.Println("Warning, there is no cli parameter named", v)
				i++
			}
		}
	}

	switch strings.ToLower(cli_params.Addon) {

	case "lang":
		compiler.Run(cli_params)
	case "compile":
		oat.Compile(cli_params)
	case "run":
		oat.Run(cli_params)
	default:
		fmt.Println("Error: cannot use omm addon", cli_params.Addon)
	}
}
