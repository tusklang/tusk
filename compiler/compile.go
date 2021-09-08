package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/initialize"
	"github.com/tusklang/tusk/varprocessor"
)

var processor = varprocessor.NewProcessor()

func Compile(prog *initialize.Program, outfile string) {

	var compiler ast.Compiler
	m := ir.NewModule() //create a new llvm module

	//set the module stuff in the compiler
	compiler.Module = m
	compiler.OS = runtime.GOOS
	compiler.ARCH = runtime.GOARCH

	compiler.StaticGlobals = make(map[string]*ir.Global)
	compiler.VarMap = make(map[string]*data.Variable)

	initDefaultOps(&compiler)
	inputDefaultTypes(&compiler)

	//initialize the operations
	var initfunc = m.NewFunc("_init", types.Void) //initialize func ran before main
	compiler.InitBlock = initfunc.NewBlock("")

	//add all the classes (files) to the type list
	for _, v := range prog.Packages {
		for k, vv := range v.Files {
			stype := types.NewStruct() //create a new structure (representing a class)

			//construct the classname
			//made up of the package name and the class' name
			var classname = vv.Name //if there isn't a package name, then it's just the class' name

			if v.Name != nil {
				classname = strings.Join(v.Name, ".") + "." + vv.Name
			}

			definiton := m.NewTypeDef(classname, stype)

			compiler.ValidTypes[classname] = definiton
			v.Files[k].StructType = stype
		}
	}

	//process all the variables

	for _, v := range prog.Packages {
		for _, vv := range v.Files {
			processor.ProcessVars(vv)
		}
	}

	j, _ := json.MarshalIndent(prog, "", "  ")
	fmt.Println(string(j))

	for _, v := range prog.Packages { //go through every package
		for _, vv := range v.Files { //go through every file in the package (class)

			//loop through all the global variables in the class
			for _, vvv := range vv.Globals {
				vvv.Value.CompileGlobal(&compiler, vv.StructType, vvv.IsStatic)
			}

		}
	}

	//declare the llvm entry function
	mfunc := m.NewFunc("main", types.Void)
	mblock := mfunc.NewBlock("")

	var (
		mfnc   *ir.Global
		exists bool
	)

	if mfnc, exists = compiler.StaticGlobals[prog.Config.Entry+"_main"]; !exists {
		//no main function
		//error
	}

	_ = mfnc

	compiler.InitBlock.NewRet(nil) //append a `return void` to the init function

	mblock.NewCall(initfunc) //call the initialize function
	loaded := mblock.NewLoad(mfnc.ContentType, mfnc)
	mblock.NewCall(loaded)
	mblock.NewRet(nil)

	f, _ := os.Create(outfile)
	f.Write([]byte(m.String()))
}
