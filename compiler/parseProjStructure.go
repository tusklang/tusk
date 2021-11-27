package compiler

import (
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/data"
	"github.com/tusklang/tusk/parser"
)

func convertPackages(compiler *ast.Compiler, packs []*parser.Package, parent *data.Package, classlist map[*parser.File]*data.Class) {

	for _, v := range packs {
		dpack := data.NewPackage(v.Name, v.FullName(), parent)
		parent.ChildPacks[v.Name] = dpack

		for _, vv := range v.Files {
			tc := compileClass(compiler, vv, v, dpack)
			dpack.Classes[vv.Name] = tc
			classlist[vv] = tc
		}

		convertPackages(compiler, v.ChildPacks, dpack, classlist)
	}

}

func parseProjStructure(compiler *ast.Compiler, prog *parser.Program) map[*parser.File]*data.Class {
	var superpack = data.NewPackage("super", "super", nil)
	var classlist = make(map[*parser.File]*data.Class) //to store all the classes compiled, regardless of nested-ness
	convertPackages(compiler, prog.Packages, superpack, classlist)

	for _, v := range superpack.ChildPacks {
		v.RemParent() //remove the superpack as the parent
		prevars[v.PackageName] = v
	}

	return classlist
}
