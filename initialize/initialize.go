package initialize

import "github.com/tusklang/tusk/ast"

//this package is used to initialize nested functions, OOP, and other high level concepts that the llvm IR can't comprehend

func Initialize(ast []*ast.ASTNode) *File {

	var f File

	fetchGlobals(ast, &f, &f.Private, 0)

	return &f
}
