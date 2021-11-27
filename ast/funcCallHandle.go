package ast

import "github.com/tusklang/tusk/errhandle"

func funcCallHandle(items []Group, i int) ([]*ASTNode, *errhandle.TuskError) {

	var argsG = items[i].(*Block)
	argsG.BlockType = "fncallb" //set the block type to a function call block

	fc, e1 := groupsToAST(items[:i])
	args, e2 := groupsToAST([]Group{argsG})

	if e2 != nil {
		return nil, e2
	}

	return []*ASTNode{
		{
			Left:  fc,
			Right: args,
			Group: &Operation{
				OpType: "()",
				tok:    items[i].GetMTok(),
			},
		},
	}, e1
}
