package compiler

//functions are written like:
//  `function_name` sync|async [params]
//but can be written like:
//  function_name[params]
//which will automatically make it sync
//because of this file

func funcLex(lex []Lex) ([]Lex, CompileErr) {
  var nLex []Lex

  for k, v := range lex {
    nLex = append(nLex, v)

    if (v.Type != "operation" && v.Type != "?operation" && v.Type != "?open_brace" && v.Type != "id" && v.Type != "id_non_tilde") && k + 2 <= len(lex) && lex[k + 1].Name == "[" {
      return []Lex{}, makeCompilerErr("Expected a [ instead of a ( for a function call", lex[k + 2].Dir, lex[k + 2].Line)
    }

    if (v.Type != "operation" && v.Type != "?operation" && v.Type != "?open_brace" && v.Type != "id" && v.Type != "id_non_tilde") && k + 2 <= len(lex) && lex[k + 1].Name == "[" {
      //insert a "sync"
      nLex = append(nLex, Lex{
        Name: "<-",
        Exp: v.Exp,
        Line: v.Line,
        Type: "operation",
        OName: "sync",
        Dir: v.Dir,
      })
    }
  }

  return nLex, nil
}
