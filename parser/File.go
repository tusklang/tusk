package parser

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/tusklang/tusk/ast"
)

type GlobalDecl struct {

	/*
		0: public
		1: protected
		2: private
	*/
	Access int

	/*
		0: instance
		1: static
		2: link
	*/
	CRel        int
	Value       *ast.VarDecl
	Link        *ast.Link
	Func        *ast.Function
	Declaration *ir.Global
}

type File struct {
	Name        string
	Globals     []GlobalDecl
	StructType  *types.StructType
	Constructor *ast.Construct
}
