package compiler

import (
	"fmt"
	"runtime"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/initialize"
)

func Compile(prog *initialize.Program, outfile string) {

	var compiler Compiler
	m := ir.NewModule() //create a new llvm module

	//set the module stuff in the compiler
	compiler.module = m
	compiler.OS = runtime.GOOS
	compiler.ARCH = runtime.GOARCH

	for _, v := range prog.Packages { //go through every package
		for _, vv := range v.Files { //go through every file in the package (class)
			stype := types.NewStruct() //create a new structure (representing a class)
			definiton := m.NewTypeDef(vv.Name, stype)
			_ = definiton

			//loop through all the global variables in the class
			for _, vvv := range vv.Globals {
				vvv.Value.CompileGlobal(stype)
			}

		}
	}

	fmt.Println(m.String())
}
