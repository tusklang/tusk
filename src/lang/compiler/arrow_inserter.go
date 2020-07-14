package compiler

var arrow_ids = []string{ "if", "elif", "else", "while", "each", "function", "fargc" }

func insert_arrows(lex []Lex) []Lex {

  //Omm needs arrows between ) and { in almost everthing. For example, this
  //  while (true) {}
  //would become
  //  while (true) => {}
  //this function does that automatically

  var nLex []Lex

  for i := 0; i < len(lex); i++ {

    nLex = append(nLex, lex[i])

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
