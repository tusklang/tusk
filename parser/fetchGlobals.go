package parser

import (
	"github.com/tusklang/tusk/ast"
	"github.com/tusklang/tusk/errhandle"
)

func fetchGlobals(tree []*ast.ASTNode, file *File, access int, crel int /*class relation, instance is 0, static is 1, and link is 2*/) *errhandle.TuskError {

	for _, v := range tree {

		switch g := v.Group.(type) {
		case *ast.VarDecl:
			file.Globals = append(file.Globals, GlobalDecl{
				Access: access,
				CRel:   crel,
				Value:  g,
			})
		case *ast.Function:
			file.Globals = append(file.Globals, GlobalDecl{
				Access: access,
				CRel:   crel,
				Func:   g,
			})
		case *ast.Public:
			return fetchGlobals([]*ast.ASTNode{g.Declaration}, file, 0, crel)
		case *ast.Protected:
			return fetchGlobals([]*ast.ASTNode{g.Declaration}, file, 1, crel)
		case *ast.Private:
			return fetchGlobals([]*ast.ASTNode{g.Declaration}, file, 2, crel)
		case *ast.Static:
			return fetchGlobals([]*ast.ASTNode{g.Declaration}, file, access, 1)
		case *ast.Link:

			g.Access = access

			file.Globals = append(file.Globals, GlobalDecl{
				Access: access,
				CRel:   2,
				Link:   g,
			})

		case *ast.Construct:
			file.Constructor = g
		default:
			return errhandle.NewParseErrorFTok(
				"invalid global",
				"",
				g.GetMTok(),
			)
		}
	}

	return nil
}
