package ast

func termOpHandle(exp []Group, index int) ([]*ASTNode, error) {

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
