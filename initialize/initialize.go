package initialize

import (
	"encoding/json"
	"fmt"
	"os"
)

//this package is used to initialize the program
/*
	- comprehending project structure
	- taking high level tusk (curried/nested functions, classes, polymorphism, etc..) and simplifying them for the llvm ir
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

	j, _ := json.MarshalIndent(config, "", "  ")
	fmt.Println(string(j))

	return &prog
}
