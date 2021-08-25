package operations

import "github.com/tusklang/tusk/grouper"

func defaultOperationHandle(exp []grouper.Group, index int) ([]*Operation, error) {

	var (
		//get the first and second half of the expression
		first  = exp[:index]
		second = exp[index+1:]

		//get the first and second half as operations
		firstop, e1  = operationsParser(first)
		secondop, e2 = operationsParser(second)
	)

	//there was an error with the sub-operation parsing
	if e1 != nil || e2 != nil {
		if e1 != nil { //if the error is in e1, move it to e2
			e2 = e1
		}
		return nil, e2
	}

	return []*Operation{{
		Left:  firstop,
		Right: secondop,
		Group: exp[index],
	}}, nil
}
