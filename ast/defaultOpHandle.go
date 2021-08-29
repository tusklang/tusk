package ast

//handler for most operations in `groupsToAST`

func defaultOperationHandle(exp []Group, index int) ([]*ASTNode, error) {

	var (
		//get the first and second half of the expression
		first  = exp[:index]
		second = exp[index+1:]

		//get the first and second half as operations
		firstop, e1  = groupsToAST(first)
		secondop, e2 = groupsToAST(second)
	)

	//there was an error with the sub-operation parsing
	if e1 != nil || e2 != nil {
		if e1 != nil { //if the error is in e1, move it to e2
			e2 = e1
		}
		return nil, e2
	}

	return []*ASTNode{{
		Left:  firstop,
		Right: secondop,
		Group: exp[index],
	}}, nil
}
