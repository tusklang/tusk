package tokenizer

var reserved = []string{"fn", "return", "var", "if", "else", "for", "while"}

func addSpecialOps(tokens []Token) []Token {

	var fin []Token

	for _, v := range tokens {
		fin = append(fin, v)

		//check if its a reserved keyword
		//if so, but the reserved operator
		for _, vv := range reserved {
			if v.Name == vv {
				fin = append(fin, Token{
					Name: "STATEMENT-OP",
					Row:  v.Row,
					Col:  v.Col,
				})
				goto endTokLoop
			}
		}

	endTokLoop: //label to end the token loop
	}

	return fin
}
