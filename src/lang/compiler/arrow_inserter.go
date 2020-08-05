package compiler

var arrow_ids = []string{ "if", "elif", "while", "each", "function" }

func insert_arrows(lex []Lex) []Lex {

  //Omm needs arrows between ) and { in almost everthing. For example, this
  //  while (true) {}
  //would become
  //  while (true) => {}
  //this function does that automatically

  var nLex []Lex

  for i := 0; i < len(lex); i++ {

    nLex = append(nLex, lex[i])

    //insert the function arrow (func_name[] or func_name() because func_name <- [])
    {
      k := i
      v := lex[k]

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

    if func() bool { //if the current token is an arrow id
      for _, id := range arrow_ids {
        if id == lex[i].Name {
          return true
        }
      }
      return false
    }() {

      var bracetype = ""
      var closebracetype = ""
      var braceCnt = 0 //how many braces there are (could be ( or {)

      bracetype = lex[i + 1].Name

      if bracetype == "(" {
        closebracetype = ")"
      } else if bracetype == "{" {
        closebracetype = "}"
      }

      //insert a => after the () expression
      for i++; i < len(lex); i++ {
        if lex[i].Name == bracetype {
          braceCnt++
        }
        if lex[i].Name == closebracetype {
          braceCnt--
        }

        nLex = append(nLex, lex[i])

        if braceCnt == 0 {
          nLex = append(nLex, Lex{
            Name: "=>",
            Exp: lex[i].Exp,
            Line: lex[i].Line,
            Type: "operation",
            OName: "=>",
            Dir: lex[i].Dir,
          })
          break
        }

      }
    }
  }

  return nLex
}
