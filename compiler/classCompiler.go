package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/initialize"
)

func compileClass(compiler *ast.Compiler, f *initialize.File, ipack *initialize.Package, pack *data.Package) *data.Class {
	stype := types.NewStruct() //create a new structure (representing a class)

	tc := data.NewClass(f.Name, stype, pack) //create the class in tusk

	//init the instance and static maps
	tc.Instance = make(map[string]*data.Variable)
	tc.Static = make(map[string]*data.Variable)

	f.StructType = stype

	//define the type in llvm
	d := compiler.Module.NewTypeDef("tuskclass."+ipack.FullName()+f.Name, stype)

	_ = d

	//define the function to create a new instance of this class
	// initf := compiler.Module.NewFunc("tuskclass.new."+ipack.FullName()+f.Name, d)
	// tc.Construct = initf.NewBlock("")
	tc.Construct = ir.NewBlock("")

	pack.AddClass(f.Name, tc)

	return tc
}
