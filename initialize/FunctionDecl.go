package initialize

import "github.com/tusklang/tusk/ast"

type FunctionDecl struct {
	Function *ast.Function
	Capture  []string
}
