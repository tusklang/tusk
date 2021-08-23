package tokenizer

var statements = []string{"fn", "return", "var", "if", "else", "for", "while"}
var bodykeys = []string{"if", "fn", "for", "while"}

func addSpecialOps(tokens []Token) []Token {

	var fin []Token

	for i := 0; i < len(tokens); i++ {
		fin = append(fin, tokens[i])

		//check if its a statements keyword
		//if so, but the statement operator
		for _, vv := range statements {
			if tokens[i].Type == vv {
				fin = append(fin, Token{
					Type: "STATEMENT-OP",
					Row:  tokens[i].Row,
					Col:  tokens[i].Col,
				})
				break
			}
		}

		var insertedFNOP bool

		//between the fn name and the fn name
		if tokens[i].Type == "fn" && tokens[i+1].Type == "varname" {
			fin = append(fin, tokens[i+1], Token{
				Type: "FN-OP",
				Row:  tokens[i].Row,
				Col:  tokens[i].Col,
			})
			insertedFNOP = true
		}

		//body operation is inserted before curly brace blocks
		//but only if the curly brace block is for a special block
		//e.g. if () body-op {}, fn () body-op {}, for () body-op {}
		for _, vv := range bodykeys {
			if tokens[i].Name == vv {

				//if we inserted a FN-OP before, we need to skip the function name
				if insertedFNOP {
					i++
				}

				if tokens[i+1].Type != "(" {
					//if the next token is not an opening parenthesis, then it isn't a body op
					break
				}

				pCnt := 0

				//find the start of the body
				for i++; i < len(tokens); i++ {
					if tokens[i].Name == "(" {
						pCnt++
					}
					if tokens[i].Name == ")" {
						pCnt--
					}

					fin = append(fin, tokens[i])

					if pCnt == 0 {
						break
					}
				}

				fin = append(fin, Token{
					Type: "BODY-OP",
					Row:  tokens[i-1].Row,
					Col:  tokens[i-1].Col,
				})

				break
			}
		}

		//tusk needs an operator in between a function and the arguments
		//f() -> f FUNCTION-CALL ()
		//but we also need to account for one-off calls
		/*
			fn() {

			}() ->
			fn() {

			} FUNCTION-CALL ()
		*/

		//case 1
		if tokens[i].Type == "varname" && (i+1 < len(tokens) && tokens[i+1].Name == "(") {
			fin = append(fin, Token{
				Type: "FUNCTION-CALL",
				Row:  tokens[i].Row,
				Col:  tokens[i].Col,
			})
		}

	}

	return fin
}
