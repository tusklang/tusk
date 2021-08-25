package ast

import (
	"github.com/tusklang/tusk/grouper"
)

//used to group operations into a tree
/*
pub fn int main() {
	var a = 1 + 3 + 4;
}
becomes:

pub fn main
-> =
	-> var a
	-> +
		-> 1
		-> +
			-> 3
			-> 4
*/

type operation struct {
	Left  *operation    //left side operand
	Right *operation    //right side operand
	Group grouper.Group //operator group/token
}

func defaultOperationHandle(exp []grouper.Group, index int) (*operation, error) {

	var (
		//get the first and second half of the expression
		first  = exp[:index]
		second = exp[index+1:]

		//get the first and second half as operations
		firstop, e1  = OperationsParser(first)
		secondop, e2 = OperationsParser(second)
	)

	//there was an error with the sub-operation parsing
	if e1 != nil || e2 != nil {
		if e1 != nil { //if the error is in e1, move it to e2
			e2 = e1
		}
		return nil, e2
	}

	return &operation{
		Left:  firstop,
		Right: secondop,
		Group: exp[index],
	}, nil
}

func OperationsParser(items []grouper.Group) (*operation, error) {

	var operations = []map[string]func(exp []grouper.Group, index int) (*operation, error){
		{
			"terminator": defaultOperationHandle,
		},
		{
			"=": defaultOperationHandle,
		},
		{
			":": defaultOperationHandle,
		},
		{
			"+": defaultOperationHandle,
			"-": defaultOperationHandle,
		},
		{
			"*": defaultOperationHandle,
			"/": defaultOperationHandle,
		},
		{
			"**": defaultOperationHandle,
		},
		//lower on this list means greater precedence
	}

	//go through all the operation groups
	for _, v := range operations {

		//go through all the items
		for i := 0; i < len(items); i++ {

			for k, vv := range v {

				switch g := items[i].(type) {
				case *grouper.Default:
					if g.Token.Type == k {
						return vv(items, i)
					}
				}
			}

		}
	}

	if len(items) != 1 {
		//only occurs when operator doesn't have two sides (!, ++, --, etc)
		return &operation{}, nil
	}

	//it must be a single, since there is no operation
	return &operation{
		Group: items[0],
	}, nil
}
