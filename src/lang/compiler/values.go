package compiler

import "unicode"
import . "lang/types"

func valueActions(item Item) (Action, CompileErr) {

  switch item.Type {

    case "{":

      acts, e := actionizer(
        makeOperations(
          item.Group,
        ),
      )

      if e != nil {
        return Action{}, e
      }

      return Action{
        Type: "{",
        ExpAct: acts,
        File: item.File,
        Line: item.Line,
      }, nil
    case "(":

      acts, e := actionizer(
        makeOperations(
          item.Group,
        ),
      )

      if e != nil {
        return Action{}, e
      }

      return Action{
        Type: "(",
        ExpAct: acts,
        File: item.File,
        Line: item.Line,
      }, nil
    case "[:":

      hash := map[string][]Action{}

      for _, v := range item.Group {
        oper := makeOperations([][]Item{ v })[0]

        //give errors
        if oper.Type != ":" {
          return Action{}, makeCompilerErr("Expected a ':' for a hash key", item.File, oper.Line)
        }
        if oper.Left.Type != "none" {
          return Action{}, makeCompilerErr("Only basic types can be used as hash indexes", item.File, oper.Line)
        }
        /////////////

        key := (*oper.Left).Item.Token.Name

        //if it is a string or rune, remove the quotes
        if key[0] == '\'' || key[0] == '"' || key[0] == '`' {
          key = key[1:len(key) - 1]
        }
        if key[0] == '$' { //if it is a variable, remove the $
          key = key[1:]
        }

        value, e := actionizer([]Operation{ *oper.Right })

        if e != nil {
          return Action{}, e
        }

        if len(value) == 0 {
          return Action{}, makeCompilerErr("Expected some value as after ':'", item.File, oper.Line)
        }

        hash[key] = value
      }

      return Action{
        Type: "hash",
        Hash: hash,
        File: item.File,
        Line: item.Line,
      }, nil
    case "[":

      var arr [][]Action

      for _, v := range item.Group {
        oper := makeOperations([][]Item{ v })[0]
        value, e := actionizer([]Operation{ oper })

        if e != nil {
          return Action{}, e
        }

        if len(value) == 0 {
          return Action{}, makeCompilerErr("Each entry in the array must have a value", item.File, item.Line)
        }

        arr = append(arr, value)
      }

      return Action{
        Type: "array",
        Array: arr,
        File: item.File,
        Line: item.Line,
      }, nil
    case "expression value":

      var val = item.Token.Name

      if val[0] == '"' || val[0] == '`' { //detect string
        var str = OmmString{}
        str.FromGoType(val[1:len(val) - 1])
        return Action{
          Type: "string",
          Value: str,
          File: item.File,
          Line: item.Line,
        }, nil
      } else if val[0] == '\'' { //detect a rune
        var oRune = OmmRune{}

        qrem := val[1:len(val) - 1] //remove quotes

        if len(qrem) != 1 {
          return Action{}, makeCompilerErr("Runes must be one character long", item.File, item.Line)
        }

        oRune.FromGoType([]rune(qrem)[0])
        return Action{
          Type: "rune",
          Value: oRune,
          File: item.File,
          Line: item.Line,
        }, nil
      } else if val == "true" || val == "false" { //detect a bool
        var boolean = OmmBool{}
        boolean.FromGoType(val == "true" /* convert to a boolean */)
        return Action{
          Type: "bool",
          Value: boolean,
          File: item.File,
          Line: item.Line,
        }, nil
      } else if val == "undef" { //detect a falsey value
        var undef OmmUndef
        return Action{
          Type: "undef",
          Value: undef,
          File: item.File,
          Line: item.Line,
        }, nil
      } else if unicode.IsDigit(rune(val[0])) || val[0] == '.' || val[0] == '+' || val[0] == '-' { //detect a number
        var number = OmmNumber{}
        number.FromString(val)
        return Action{
          Type: "number",
          Value: number,
          File: item.File,
          Line: item.Line,
        }, nil
      } else if val[0] == '$' { //detect a variable
        return Action{
          Type: "variable",
          Name: val,
          File: item.File,
          Line: item.Line,
        }, nil
      } else { //detect nothing, which throws an error
        return Action{}, makeCompilerErr(val + " is not a value", item.File, item.Line)
      }

  }

  return Action{}, nil
}
