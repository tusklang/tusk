package compiler

import "path"
import "strings"
import "unicode"
import "encoding/json"
import "regexp"

type Lex struct {
  Name   string
  Exp    string
  Line   uint64
  Type   string
  OName  string
  Dir    string
}

//length of each error expression
const EXPRESSION_LEN = 30

func lexer(file, dirname, filename string) []Lex {

  var keywords []map[string]string

  _ = json.Unmarshal(keywordJSON, &keywords)

  var lex []Lex
  curExp := ""
  var line uint64 = 1

  for i := 0; i < len(file); i++ {

    for ;len(file[i:]) != 0 && unicode.IsSpace(rune(file[i:][0])); i++ {}

    if len(strings.TrimSpace(file[i:])) == 0 {
      break
    }

    //detect a comment
    //single line comments are written as ;comment
    //like in assembly
    if strings.TrimSpace(file[i:])[0] == ';' {

      var end = strings.Index(file[i:], "\n")

      //if there is no newline after the comment, just break the loop
      if end == -1 {
        break
      }

      i+=end
      line++
      continue
    }

    for ;unicode.IsSpace(rune(file[i])); {
      continue
    }

    //while the curExp is too long, pop the front
    for ;len(curExp) > EXPRESSION_LEN; curExp = curExp[1:] {}

    for _, v := range keywords {

      if testkey(v, file, i) {
        if v["name"] == "newlineN" { //if it is a newline, increment the line
          line++
        }

        curExp+=v["remove"] + " "
        i+=len(v["remove"])

        lex = append(lex, Lex{
          Name: v["name"],
          Exp: curExp,
          Line: line,
          Type: v["type"],
          OName: v["remove"],
          Dir: path.Join(dirname, filename),
        })
        goto contOuter
      }

    }

    if strings.TrimSpace(file)[i:][0] == '"' || strings.TrimSpace(file)[i:][0] == '\'' || strings.TrimSpace(file)[i:][0] == '`' { //detect a string
      qType := strings.TrimSpace(file)[i:][0]
      value := ""
      escaped := false

      for i++; i < len(file); i++ {
        value+=string(strings.TrimSpace(file)[i:][0])

        if !escaped && strings.TrimSpace(file)[i:][0] == '\\' {
          escaped = true
          continue
        }

        if !escaped && strings.TrimSpace(file)[i:][0] == qType {
          break
        }
        escaped = false
      }

      curExp+=value + " "
      line+=uint64(strings.Count(value, "\n"))
      i++
      lex = append(lex, Lex{
        Name: "\"" + value[:len(value) - 1] + "\"",
        Exp: curExp,
        Line: line,
        Type: "expression value",
        OName: value,
        Dir: path.Join(dirname, filename),
      })
      goto contOuter
    } else if unicode.IsDigit(rune(strings.TrimSpace(file)[i:][0])) || strings.TrimSpace(file)[i:][0] == '+' || strings.TrimSpace(file)[i:][0] == '-' || strings.TrimSpace(file)[i:][0] == '.' {

      var positive = true

      if strings.TrimSpace(file)[i:][0] == '+' {
        positive = true
        i++
      } else if strings.TrimSpace(file)[i:][0] == '-' {
        positive = false
        i++
      }

      num := ""

      for ;unicode.IsDigit(rune(strings.TrimSpace(file)[i:][0])) || strings.TrimSpace(file)[i:][0] == '.'; i++ {
        num+=string(strings.TrimSpace(file)[i:][0])
      }

      if !positive {
        num = "-" + num
      }

      curExp+=num + " "

      lex = append(lex, Lex{
        Name: num,
        Exp: curExp,
        Line: line,
        Type: "expression value",
        OName: num,
        Dir: path.Join(dirname, filename),
      })
    } else {

      if unicode.IsSpace(rune(file[i:][0])) {
        continue
      }

      variable := ""

      for o := i; o < len(file); o++ {
        for _, v := range keywords {
          if testkey(v, file, o) || unicode.IsSpace(rune(file[o])) || file[o] == ';' /* it is a comment */ {
            goto break_var_loop
          }
        }

        variable+=string(file[o])
        i++
      }
      break_var_loop:

      curExp+=variable + " "
      variable = strings.TrimSpace(variable)
      i--

      lex = append(lex, Lex{
        Name: "$" + variable,
        Exp: curExp,
        Line: line,
        Type: "expression value",
        OName: variable,
        Dir: path.Join(dirname, filename),
      })
    }

    contOuter: //continue the outer loop
  }

  //filter out newlines
  var newLex []Lex
  for _, v := range lex {
    if v.Type == "newlineN" {
      continue
    }
    if v.Name == "$false" || v.Name == "$true" || v.Name == "$undef" || v.Name == "$null" { //account for true, false, undef, and null values
      v.Name = v.Name[1:]
      newLex = append(newLex, v)
      continue
    }
    newLex = append(newLex, v)
  }
  lex = newLex

  lex = term_inserter(tilde_inserter(lex))

  return lex
}

func testkey(keyword map[string]string, file string, i int) bool {
  re, _ := regexp.Compile("^(" + keyword["pattern"] + ")")
  return re.MatchString(strings.TrimSpace(file[i:]))
}
