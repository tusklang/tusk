package ast

import "github.com/tusklang/tusk/errhandle"

func arrIndexHandle(items []Group, i int) ([]*ASTNode, *errhandle.TuskError) {
	var idx = items[i].(*Array)
	idx.useAsIndex = true

	fc, e1 := groupsToAST(items[:i])
	args, e2 := groupsToAST([]Group{idx})

	if e2 != nil {
		return nil, e2
	}

	return []*ASTNode{
		{
			Left:  fc,
			Right: args,
			Group: &Operation{
				OpType: "[]",
				tok:    items[i].GetMTok(),
			},
		},
	}, e1
}
