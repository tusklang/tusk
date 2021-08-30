package ast

type ASTNode struct {
	Left   []*ASTNode //left side operand
	Right  []*ASTNode //right side operand
	Group  Group      //operator group/token
	parent *ASTNode   //parent node
}

func (a *ASTNode) Parent() *ASTNode {
	return a.parent
}
