package ast

import (
	"github.com/tusklang/tusk/errhandle"
	"github.com/tusklang/tusk/tokenizer"
)

//package to group major areas of lex together
//such as function headers, if statement conditions, codeblocks, variable declarations, etc..

/*
fn main() {
	if (true) {
		var a = 3 + 1;
	}
}

becomes

{
	gtype(): fn-head
	name: main
	params: []
},
{
	gtype(): block
	blocktype: {
	sub: [
		{
			gtype(): if
			condition: [
				{
					gtype(): value
					data: true
				}
			]
		},
		{
			gtype(): block
			blocktype: {
			sub: [
				{
					gtype(): var
					name: a
				},
				{
					gtype(): operator
					op: =
				},
				{
					gtype(): value
					data: 3
				},
				{
					gtype(): operator
					op: =
				},
				{
					gtype(): value
					data: 1
				},
				{
					gtype(): terminator
				}
			]
		}
	]
}
*/

func testStopAt(token tokenizer.Token, sa []string) bool {
	for _, v := range sa {
		if token.Name == v {
			return false
		}
	}
	return true
}

//function for extra customizability with grouping
func groupSpecific(tokens []tokenizer.Token, startAt *int, stopAt []string, maxlen int) ([]Group, *errhandle.TuskError) {
	var fin []Group
	ostart := *startAt

	for ; *startAt < len(tokens) && testStopAt(tokens[*startAt], stopAt) && (maxlen < 0 || *startAt < ostart+maxlen); *startAt++ {

		var gr Group //the group to append

		switch tokens[*startAt].Type {
		case "fn":
			gr = &Function{}
		case "return":
			gr = &Return{}
		case "{":
			gr = &Block{}
		case "(":
			gr = &Block{}
		case "[":
			gr = &Array{}
		case "if":
			gr = &IfStatement{}
		case "while":
			gr = &WhileStatement{}
		case "pub":
			gr = &Public{}
		case "prt":
			gr = &Protected{}
		case "prv":
			gr = &Private{}
		case "stat":
			gr = &Static{}
		case "link":
			gr = &Link{}
		case "pure":
			gr = &Pure{}
		case "var":
			gr = &VarDecl{}
		case "terminator":
			fallthrough
		case "operation":
			gr = &Operation{}
		case "float":
			fallthrough
		case "bool":
			fallthrough
		case "int":
			gr = &DataValue{}
		case "string":
			gr = &String{}
		case "null":
			gr = &NullValue{}
		case "varname":
			gr = &VarRef{}
		case "construct":
			gr = &Construct{}
		case "this":
			gr = &This{}
		default:
			//error
			//the token given isn't recognized by tusk
		}

		e := gr.Parse(tokens, startAt, stopAt)

		if e != nil {
			return nil, e
		}

		fin = append(fin, gr)

	}

	return fin, nil
}

//function used as shorthand for `groupSpecific` when some params aren't required
func grouper(tokens []tokenizer.Token) ([]Group, *errhandle.TuskError) {
	tmp := 0
	return groupSpecific(tokens, &tmp, nil, -1)
}
