package grouper

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

func Grouper(tokens []tokenizer.Token) []Group {
	var fin []Group

	for i := 0; i < len(tokens); i++ {

		var gr Group //the group to append

		switch tokens[i].Type {
		case "fn":
			gr = &FunctionHeader{}
		case "{":
			fallthrough
		case "(":
			gr = &Block{}
		case "if":
			gr = &IfStatement{}
		case "while":
			gr = &WhileStatement{}
		default:
			gr = &Default{}
		}

		_ = gr.Parse(tokens, &i)

		fin = append(fin, gr)
	}

	return fin
}
