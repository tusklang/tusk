package compiler

//functions are written like:
//  `function_name` sync|async (params)
//but can be written like:
//  function_name(params)
//which will automatically make it sync
//because of this file

func funcLex(lex []Lex) []Lex {
  var nLex []Lex

  for k, v := range lex {
    nLex = append(nLex, v)

    if (v.Name[0] == '$' || v.Name == ")" || v.Name == "}") && k + 2 <= len(lex) && lex[k + 1].Name == "(" {
      //insert a "sync"
      nLex = append(nLex, Lex{
        Name: "sync",
        Exp: v.Exp,
        Line: v.Line,
        Type: "operation",
        OName: "sync",
        Dir: v.Dir,
      })
    }
  }

  return nLex
}
