package ast

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

func groupsToAST(items []Group) ([]*ASTNode, error) {

	var opList = []map[string]func(exp []Group, index int) ([]*ASTNode, error){
		{
			";": termOpHandle,
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
	for _, v := range opList {

		//go through all the items
		for i := 0; i < len(items); i++ {

			for k, vv := range v {

				switch g := items[i].(type) {
				case *Operation:
					if g.Token.Name == k {
						return vv(items, i)
					}
				}
			}

		}
	}

	if len(items) != 1 {
		//only occurs when operator doesn't have two sides (!, ++, --, etc)
		return nil, nil
	}

	//it must be a single, since there is no operation
	return []*ASTNode{{
		Group: items[0],
	}}, nil
}
