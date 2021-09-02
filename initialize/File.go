package initialize

import "github.com/tusklang/tusk/ast"

type AccessLevel struct {
	Static, Instance []*ast.ASTNode
}

type File struct {
	Name      string
	Public    AccessLevel
	Protected AccessLevel
	Private   AccessLevel
}
