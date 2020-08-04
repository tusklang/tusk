package compiler

//functions are written like:
//  `function_name` sync|async [params]
//but can be written like:
//  function_name[params]
//which will automatically make it sync
//because of this file

func funcLex(lex []Lex) ([]Lex, CompileErr) {
  var nLex []Lex

  for k := 0; k < len(lex); k++ {
    v := lex[k]
    nLex = append(nLex, v)

    if (v.Type != "operation" && v.Type != "?operation" && v.Type != "?open_brace" && v.Type != "id" && v.Type != "id_non_tilde") && k + 2 <= len(lex) && lex[k + 1].Name == "(" { //if the dev used a ( for a function call instead of a [

      //insert a "sync"
      nLex = append(nLex, Lex{
        Name: "<-",
        Exp: v.Exp,
        Line: v.Line,
        Type: "operation",
        OName: "sync",
        Dir: v.Dir,
      })

      //insert a "["
      nLex = append(nLex, Lex{
        Name: "[",
        Exp: v.Exp,
        Line: v.Line,
        Type: "?open_brace",
        OName: "[",
        Dir: v.Dir,
      })

      pCnt := 1

      for k+=2; k < len(lex); k++ {

        if lex[k].Name == "(" {
          pCnt++
        }
        if lex[k].Name == ")" {
          pCnt--
        }

        if pCnt == 0 {
          break
        }

        nLex = append(nLex, lex[k])

      }

      //insert a "]"
      nLex = append(nLex, Lex{
        Name: "]",
        Exp: v.Exp,
        Line: v.Line,
        Type: "?close_brace",
        OName: "[",
        Dir: v.Dir,
      })
      
      continue
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
