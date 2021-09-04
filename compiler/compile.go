package compiler

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/initialize"
)

//all of the basic types
func inputDefaultTypes(compiler *ast.Compiler) {

	compiler.ValidTypes = make(map[string]types.Type)

	compiler.ValidTypes["i64"] = types.I64
	compiler.ValidTypes["i32"] = types.I32
	compiler.ValidTypes["i16"] = types.I16
	compiler.ValidTypes["i8"] = types.I8

	compiler.ValidTypes["f64"] = types.Double
	compiler.ValidTypes["f32"] = types.Float
}

func Compile(prog *initialize.Program, outfile string) {

	var compiler ast.Compiler
	m := ir.NewModule() //create a new llvm module

	//set the module stuff in the compiler
	compiler.Module = m
	compiler.OS = runtime.GOOS
	compiler.ARCH = runtime.GOARCH

	inputDefaultTypes(&compiler)

	for _, v := range prog.Packages { //go through every package
		for _, vv := range v.Files { //go through every file in the package (class)
			stype := types.NewStruct() //create a new structure (representing a class)

			//construct the classname
			//made up of the package name and the class' name
			var classname = vv.Name //if there isn't a package name, then it's just the class' name

			if v.Name != nil {
				classname = strings.Join(v.Name, ".") + "." + vv.Name
			}

			definiton := m.NewTypeDef(classname, stype)

			compiler.ValidTypes[classname] = definiton

			//loop through all the global variables in the class
			for _, vvv := range vv.Globals {
				vvv.Value.CompileGlobal(&compiler, stype)
			}

		}
	}

	fmt.Println(m.String())
}
