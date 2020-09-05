package compiler

//insert tildes after ids
func tilde_inserter(lex []Lex) []Lex {

	var nLex []Lex

	for _, v := range lex {
		nLex = append(nLex, v)

		if v.Type == "id" {
			nLex = append(nLex, Lex{
				Name:  "~",
				Exp:   v.Exp,
				Line:  v.Line,
				Type:  "operation",
				OName: "~",
				Dir:   v.Dir,
			})
		}
	}

	return nLex
}
