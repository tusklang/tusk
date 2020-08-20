package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/omm-lang/omm/oat"

	. "github.com/omm-lang/omm/lang/types"

	"github.com/omm-lang/omm/lang/compiler"
)

//omm addons
//omm language (compile into go slices and structs)

////////////

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

	if len(args) <= 2 {
		fmt.Println("Error, no input file was given")
		os.Exit(1)
	}

	defaults(&cli_params, args[2])

	cli_params.Directory = args[1]
	cli_params.Name = args[2]

	__dirname, _ := os.Getwd()

	cli_params.OmmDirname = __dirname

	//set the working directory
	os.Chdir(args[1])

	for i := 2; i < len(args); i++ {

		v := args[i]

		if strings.HasPrefix(v, "--") {

			switch strings.ToUpper(v) {

			case "--version":
				fmt.Printf("Omm Beta %d.%d.%d", OMM_MAJOR, OMM_MINOR, OMM_BUG)
				os.Exit(0)
			default:
				cli_params.Addon = v[2:]
			}

		} else if strings.HasPrefix(v, "-") {

			switch strings.ToUpper(v[1:]) {

			case "V":
				fmt.Printf("Omm Beta %d.%d.%d", OMM_MAJOR, OMM_MINOR, OMM_BUG)
				os.Exit(0)
			case "C":
				cli_params.Addon = "compile"
			case "R":
				cli_params.Addon = "run"
			case "PREC":
				temp_prec, _ := strconv.ParseUint(args[i+1], 10, 64)
				cli_params.Prec = temp_prec
				i += 2
			case "O":
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
