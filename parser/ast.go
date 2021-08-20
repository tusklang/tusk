package parser

import "fmt"

type astNode struct {
	Name  string
	Left  []*astNode
	Right []*astNode
}

var asterr error

func convertOperation(op *operation) []*astNode {

	switch op.Item.Token.Type {

	case "terminator":

		l, r := convertOperation(op.Left), convertOperation(op.Right)

		l = append(l, r...)

		return l

	//keywords
	case "fn":
		fallthrough
	case "return":
		fallthrough
	case "var":
		fallthrough
	case "if":
		fallthrough
	case "else":
		fallthrough
	case "for":
		fallthrough
	case "while":
		return []*astNode{
			{
				Name: op.Item.Token.Type,
			},
		}
	//////////

	//braces/scopes
	case "(":
		fallthrough
	case "{":
		fallthrough
	case "[":

		group := op.Item.Sub
		ast := genAST(group)

		return []*astNode{
			{
				Name: op.Item.Token.Type,
				Left: ast,
			},
		}
	///////////////

	//special operators
	case "STATEMENT-OP":

		//get which statement it is
		stat := convertOperation(op.Left)

		fmt.Println(op.Left)

		if len(stat) != 1 {
			//error
		}

		acts := convertOperation(op.Right)

		if len(acts) != 1 {
			//error
		}

		if acts[0].Name == "BODY-OP" {
			return []*astNode{
				{
					Name:  stat[0].Name,
					Left:  acts[0].Left,
					Right: acts[0].Right,
				},
			}
		}

		return []*astNode{
			{
				Name: stat[0].Name,
				Left: acts,
			},
		}

	//these are special operators, but follows the same pattern as the normal operators
	case "BODY-OP":
		fallthrough

	case "FUNCTION-CALL":
		fallthrough

	///////////////////

	//operators
	case ":":
		fallthrough
	case "=":
		fallthrough
	case "+":
		fallthrough
	case "-":
		fallthrough
	case "*":
		fallthrough
	case "/":
		fallthrough
	case "**":

		return []*astNode{
			{
				Name:  op.Item.Token.Type,
				Left:  convertOperation(op.Left),
				Right: convertOperation(op.Right),
			},
		}

	///////////

	case "float":
		fallthrough
	case "int":
		fallthrough
	case "varname":

		var tname string = op.Item.Token.Type

		return []*astNode{ //return the final representation of the value
			{
				Name: tname,
				Left: []*astNode{
					{
						Name: op.Item.Token.Name,
					},
				},
			},
		}

	}

	return nil
}

func genAST(groups []gItem) []*astNode {

	operation, e := parseOperations(groups)

	if e != nil {
		asterr = e
		return nil
	}

	return convertOperation(operation)
}
