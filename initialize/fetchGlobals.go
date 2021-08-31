package initialize

import "github.com/tusklang/tusk/ast"

func fetchGlobals(tree []*ast.ASTNode, file *File, output *[]*ast.ASTNode) {

	for _, v := range tree {
		switch g := v.Group.(type) {
		case *ast.Operation:
			fetchGlobals(v.Left, file, output)
			//we don't check the right side because the only operator we're looking for is the '='
			//and the `var x` is always on the left
		case *ast.Function:
			*output = append(*output, v)
		case *ast.VarDecl:
			*output = append(*output, v)
		case *ast.Public:
			fetchGlobals([]*ast.ASTNode{g.Declaration}, file, &file.Public)
		}
	}
}
