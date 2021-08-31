package initialize

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//this package is used to initialize the program
/*
	- comprehending project structure
	- taking high level tusk (curried/nested functions, classes, polymorphism, etc..) and simplifying them for the llvm ir
	- dependency managment
*/

func parsePackage(dir string, pkg *Package, prog *Program) {

	fsinfo, _ := ioutil.ReadDir(dir)

	for _, v := range fsinfo {
		if v.IsDir() {
			//a new package

			var spkg Package

			//joined path (of the parent directories and current one)
			jpth := path.Join(dir, v.Name())

			//because the name is an array (see `Package.go`) we want to get the package names of all the parents
			spkg.Name = append(spkg.Name, strings.Split(jpth, "/")...)

			parsePackage(jpth, &spkg, prog)
			continue
		}

		//only append a new class if it's a tusk file
		if filepath.Ext(v.Name()) != ".tusk" {
			continue
		}

		//a new class in the package
		pkg.Files = append(pkg.Files, &File{
			Name: strings.TrimSuffix(v.Name(), ".tusk"),
		})
	}

	prog.Packages = append(prog.Packages, pkg)
}

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

	var startpkg Package

	parsePackage(".", &startpkg, &prog)

	j, _ := json.MarshalIndent(prog, "", "  ")
	fmt.Println(string(j))

	return &prog
}
