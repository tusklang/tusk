package compiler

import (
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/initialize"
)

func compileClass(compiler *ast.Compiler, f *initialize.File, ipack *initialize.Package, pack *data.Package) *data.Class {
	stype := types.NewStruct() //create a new structure (representing a class)

	tc := data.NewClass(ipack.FullName()+f.Name, stype, pack) //create the class in tusk
	tc.CName = f.Name

	//init the instance and static maps
	tc.Instance = make(map[string]*data.ClassField)
	tc.Static = make(map[string]*data.ClassField)

	f.StructType = stype

	//define the type in llvm
	d := compiler.Module.NewTypeDef("tuskclass."+ipack.FullName()+f.Name, stype)

	//define the function to create a new instance of this class
	initf := compiler.Module.NewFunc("tuskclass.new."+ipack.FullName()+f.Name, types.NewPointer(d))
	tc.Construct = data.NewFunc(initf, tc)
	tc.Construct.ActiveBlock = initf.NewBlock("")

	pack.AddClass(f.Name, tc)

	return tc
}
