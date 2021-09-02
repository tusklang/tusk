package ast

import (
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

//helper function to compare the max group size to return and the current index
func cmpMaxGroup(i, maxGroup, originalStartPos int) bool {
	if maxGroup < 0 {
		return true //continue if -1 is given, -1 means no max
	}

	return i-originalStartPos < maxGroup //otherwise just return if i is less than the max group size
}

//function for extra customizability with grouping
func groupSpecific(tokens []tokenizer.Token, maxGroup int, startAt *int) []Group {
	var fin []Group

	originalStartPos := *startAt

	for ; *startAt < len(tokens) && cmpMaxGroup(*startAt, maxGroup, originalStartPos); *startAt++ {

		var gr Group //the group to append

		switch tokens[*startAt].Type {
		case "fn":
			gr = &Function{}
		case "{":
			fallthrough
		case "(":
			gr = &Block{}
		case "if":
			gr = &IfStatement{}
		case "while":
			gr = &WhileStatement{}
		case "pub":
			gr = &Public{}
		case "prt":
			gr = &Protected{}
		case "stat":
			gr = &Static{}
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
		case "varname":
			gr = &VarRef{}
		case "dtype":
			gr = &DataType{}
		default:
			//error
			//the token given isn't recognized by tusk
		}

		_ = gr.Parse(tokens, startAt)

		fin = append(fin, gr)
	}

	return fin
}

//function used as shorthand for `groupSpecific` when some params aren't required
func grouper(tokens []tokenizer.Token) []Group {
	tmp := 0
	return groupSpecific(tokens, -1, &tmp)
}
