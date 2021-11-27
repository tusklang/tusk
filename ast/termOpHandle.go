package ast

import "github.com/tusklang/tusk/errhandle"

//handle the terminator in `groupsToAST`

func termOpHandle(exp []Group, index int) ([]*ASTNode, *errhandle.TuskError) {

	//first use the default handler to get a left and right side
	defaultOp, e := defaultOperationHandle(exp, index)

	if e != nil {
		return nil, e
	}

	//then merge both sides into one
	dleft, dright := defaultOp[0].Left, defaultOp[0].Right
	dleft = append(dleft, dright...)

	return dleft, nil
}
