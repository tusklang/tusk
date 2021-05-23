package tokenizer

import "fmt"

var statements = []string{"fn", "return", "var", "if", "else", "for", "while"}
var bodykeys = []string{"fn", "if", "for", "while"}

func addSpecialOps(tokens []Token) []Token {

	var fin []Token

	for i := 0; i < len(tokens); i++ {
		fin = append(fin, tokens[i])

		//check if its a statements keyword
		//if so, but the statement operator
		for _, vv := range statements {
			if tokens[i].Name == vv {
				fin = append(fin, Token{
					Name: "STATEMENT-OP",
					Row:  tokens[i].Row,
					Col:  tokens[i].Col,
				})
				break
			}
		}

		//body operation is inserted before curly brace blocks
		//but only if the curly brace block is for a special block
		//e.g. if () body-op {}, fn x() body-op {}, for () body-op {}
		for _, vv := range bodykeys {
			if tokens[i].Name == vv {

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
					Name: "BODY-OP",
					Row:  tokens[i-1].Row,
					Col:  tokens[i-1].Col,
				})

				break
			}
		}

	}

	fmt.Println(fin)

	return fin
}
