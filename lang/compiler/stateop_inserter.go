package compiler

//insert the statement operators after statements
//e.g. var a --> var STATE-OP a
func tilde_inserter(lex []Lex) []Lex {

	var nLex []Lex

	for _, v := range lex {
		nLex = append(nLex, v)

		if v.Type == "id" {
			nLex = append(nLex, Lex{
				Name:  "STATE-OP",
				Exp:   v.Exp,
				Line:  v.Line,
				Type:  "operation",
				OName: "STATE-OP",
				Dir:   v.Dir,
			})
		}
	}

	return nLex
}
