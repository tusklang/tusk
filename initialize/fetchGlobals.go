package initialize

import (
	"github.com/tusklang/tusk/ast"
)

func fetchGlobals(acts []*ast.ASTNode, file *File, varOutput *[]Declaration, lvl int) {

	for _, v := range acts {
		switch g := v.Group.(type) {
		case *ast.Function:
			*varOutput = append(*varOutput, FunctionDecl{
				Function: g,
			})
		case *ast.Public:
			fetchGlobals([]*ast.ASTNode{g.Declaration}, file, &file.Public, lvl)
		}
	}
}
