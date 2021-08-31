package initialize

import "github.com/tusklang/tusk/ast"

type File struct {
	Name      string
	Public    []*ast.ASTNode
	Protected []*ast.ASTNode
	Private   []*ast.ASTNode
}
