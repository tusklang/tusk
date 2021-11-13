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

	parsePackage(".", &startpkg)

	prog.Packages = startpkg.ChildPacks

	if startpkg.Files != nil {
		//error
		//this means there are files that are not within a package, which is forbidden (as of current)
		//any entry files must be placed within a package, usually entry/<files go here>
		fmt.Println("error: found files not contained within a package")
	}

	return &prog
}
