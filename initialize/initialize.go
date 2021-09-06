package initialize

import (
	"encoding/json"
	"fmt"
	"os"
)

//this package is used to initialize the program
/*
	- comprehending project structure
	- dependency managment
	- compiling the project into llvm IR
*/

func Initialize(configFileName string) *Program {

	var prog Program

	//config is the tusk config file

	configFile, e := os.Open(configFileName)

	if e != nil {
		//error
		_ = e
	}

	var config ConfigData
	json.NewDecoder(configFile).Decode(&config)

	prog.Config = config

	var startpkg Package

	parsePackage(".", &startpkg, &prog)

	j, _ := json.MarshalIndent(prog, "", "  ")
	fmt.Println(string(j))

	return &prog
}
