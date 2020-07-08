package compiler

//insert $term automatically in the lexer
func term_inserter(lex []Lex) []Lex {

  var nLex []Lex

  for k, v := range lex {
    nLex = append(nLex, v)

    currentType := v.Type == "operation" || v.Type == "?operation"
    nextType := k + 2 <= len(lex) && (lex[k + 1].Type == "operation" || lex[k + 1].Type == "?operation")

    if v.Type == "open_brace" { //because opening braces don't need a $term after it
      continue
    }

    if v.Type[0] == '?' && k + 2 <= len(lex) && lex[k + 1].Type[0] == '?' { //detect types with a ? prefix
      continue
    }

    //if it looks like
    //  (1 + 3)
    //it would become
    //  ($term 1 + 3 $term)
    //this is to prevent that
    if (v.Type == "?open_brace" && k + 2 <= len(lex) && lex[k + 1].Type == "expression value") || (v.Type == "expression value" && k + 2 <= len(lex) && lex[k + 1].Type == "?close_brace") {
      continue
    }

    if currentType == nextType {
      nLex = append(nLex, Lex{
        Name: "$term",
        Exp: v.Exp,
        Line: v.Line,
        Type: "?none",
        OName: "$term",
        Dir: v.Dir,
      })
    }
  }

  return nLex
}
