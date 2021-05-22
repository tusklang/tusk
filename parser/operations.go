package parser

//used to group operations into a tree
/*
fn a() {
	log("hello")
}

becomes:

			STATEMENT-OP (inserted after tokenizing between fn and a())
				 |
			------------
			|		   |
		    fn       BODY-OP (inserted after tokenizing between a() and {})
					   |
				  ------------
				  |			 |
		          a()   FUNCTION-CALL (inserted after tokenizing between log and ("hello"))
				  			 |
					    ------------
						|		   |
					   log	    ("hello")

(sorry for the horrible formatting)
*/

type operation struct {
	Left  *operation
	Right *operation
	Item  gItem
}

func defaultOperationHandle(exp []gItem, index int) (*operation, error) {

	var (
		//get the first and second half of the expression
		first  = exp[:index]
		second = exp[index+1:]

		//get the first and second half as operations
		firstop, e1  = parseOperations(first)
		secondop, e2 = parseOperations(second)
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
		Item:  exp[index],
	}, nil
}

func parseOperations(items []gItem) (*operation, error) {

	var operations = []map[string]func(exp []gItem, index int) (*operation, error){
		{
			";": defaultOperationHandle,
		},
		{
			"STATEMENT-OP": defaultOperationHandle,
		},
		{
			"+": defaultOperationHandle,
		},
		{
			"*": defaultOperationHandle,
		},
		//lower on this list means greater precedence
	}

	//go through all the operation groups
	for _, v := range operations {

		//go through all the items
		for i := 0; i < len(items); i++ {

			for k, vv := range v {
				if items[i].Token.Name == k {
					//current operation in loop is the operation needed
					return vv(items, i)
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
		Item: items[0],
	}, nil
}
