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
			",": termOpHandle,
		},
		{
			"=":  defaultOperationHandle,
			"->": defaultOperationHandle,
		},
		{
			":": defaultOperationHandle,
		},
		{
			"&": defaultOperationHandle,
		},
		{
			"|": defaultOperationHandle,
		},
		{
			"^": defaultOperationHandle,
		},
		{
			"==": defaultOperationHandle,
			"!=": defaultOperationHandle,
		},
		{
			"<=": defaultOperationHandle,
			"<":  defaultOperationHandle,
			">":  defaultOperationHandle,
			">=": defaultOperationHandle,
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
		{
			".":  defaultOperationHandle,
			"[]": nil, //array index
			"()": nil, //function call
		},
		{
			"~": defaultOperationHandle,
		},
		//lower on this list means greater precedence
	}

	//go through all the operation groups
	for _, v := range opList {

		//go through all the items
		//reverse order for left to right associativity
		for i := len(items) - 1; i >= 0; i-- {

			for k, vv := range v {

				switch g := items[i].(type) {
				case *Operation:
					if g.Token.Name == k {
						return vv(items, i)
					}
				case *Array:

					if k == "[]" && g.siz != nil && g.arr == nil && g.typ == nil {

						if len(items) == 1 {
							return []*ASTNode{{
								Group: items[i],
							}}, nil
						}

						return arrIndexHandle(items, i)
					}

				case *Block:
					//for function calls

					//if the blocktype is a (
					if k == "()" && g.BlockType == "(" {

						if i-1 >= 0 {

							//if the item prior is a function, variable, or block
							/*
								a();
								fn() {}();
								(fn() {})();
							*/
							switch items[i-1].(type) {
							case *VarRef:
								return funcCallHandle(items, i)
							case *Function:
								return funcCallHandle(items, i)
							case *Block:
								return funcCallHandle(items, i)
							}

						}

					}
				}
			}

		}
	}

	var ret = make([]*ASTNode, len(items))

	for k, v := range items {
		ret[k] = &ASTNode{}
		ret[k].Group = v
	}

	//it must be a single, since there is no operation
	return ret, nil
}
