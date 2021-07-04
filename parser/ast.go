package parser

import (
	"github.com/dlclark/regexp2"
	"github.com/tusklang/tusk/tokenizer"
)

type astNode struct {
	Name  string
	Left  []*astNode
	Right []*astNode
}

var asterr error

func convertOperation(op *operation) []*astNode {

	switch op.Item.Token.Name {

	case ";":

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
				Name: op.Item.Token.Name,
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
				Name: op.Item.Token.Name,
				Left: ast,
			},
		}
	///////////////

	//special operators
	case "STATEMENT-OP":

		//get which statement it is
		stat := convertOperation(op.Left)

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
				Name:  op.Item.Token.Name,
				Left:  convertOperation(op.Left),
				Right: convertOperation(op.Right),
			},
		}

	///////////

	default:

		var (
			reFloat = regexp2.MustCompile(tokenizer.FloatPat, 0) //float regex
			reInt   = regexp2.MustCompile(tokenizer.IntPat, 0)   //int regex
		)

		var tname string

		if b, _ := reFloat.MatchString(op.Item.Token.Name); b {
			tname = "float"
		} else if b, _ := reInt.MatchString(op.Item.Token.Name); b {
			tname = "integer"
		} else if tokenizer.IsVariable(op.Item.Token) { //test if its a variable
			tname = "variable"
		}

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

}

func genAST(groups []gItem) []*astNode {

	operation, e := parseOperations(groups)

	if e != nil {
		asterr = e
		return nil
	}

	return convertOperation(operation)
}
