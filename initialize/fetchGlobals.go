package initialize

import "github.com/tusklang/tusk/ast"

func fetchGlobals(tree []*ast.ASTNode, file *File, accessLevel *AccessLevel, output *[]*ast.ASTNode) {

	for _, v := range tree {
		switch g := v.Group.(type) {
		case *ast.Function:
			*output = append(*output, v)
		case *ast.VarDecl:
			*output = append(*output, v)
		case *ast.Public:
			fetchGlobals([]*ast.ASTNode{g.Declaration}, file, &file.Public, &file.Public.Instance)
		case *ast.Protected:
			fetchGlobals([]*ast.ASTNode{g.Declaration}, file, &file.Protected, &file.Protected.Instance)
		case *ast.Static:
			fetchGlobals([]*ast.ASTNode{g.Declaration}, file, accessLevel, &accessLevel.Static)
		}
	}
}
