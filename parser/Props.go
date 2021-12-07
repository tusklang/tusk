package parser

func IsStatic(g GlobalDecl) bool {
	return g.CRel != 0
}

func IsPure(g GlobalDecl) bool {
	return g.CRel == 3
}
