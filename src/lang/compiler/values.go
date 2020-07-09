package compiler

import "unicode"
import . "lang/types"

func valueActions(item Item, dir string) []Action {

  switch item.Type {

    case "{":

      if len(item.Group) == 0 {
        compilerErr("{} is not allowed, try using {_} instead", dir, item.Line)
      }

      actionized := actionizer(
        makeOperations(
          item.Group,
        ),
        dir,
      )

      return actionized
    case "(":

      if len(item.Group) == 0 {
        compilerErr("() is not allowed, try using (_) instead", dir, item.Line)
      }

      return actionizer(
        makeOperations(
          item.Group,
        ),
        dir,
      )
    case "[:":

      hash := map[string]Action{}

      for _, v := range item.Group {
        oper := makeOperations([][]Item{ v })[0]

        //give errors
        if oper.Type != ":" {
          compilerErr("Expected a ':' for a hash key", dir, oper.Line)
        }
        if oper.Left.Type != "none" {
          compilerErr("Only basic types can be used as hash indexes", dir, oper.Line)
        }
        /////////////

        key := (*oper.Left).Item.Token.Name
        value := actionizer([]Operation{ *oper.Right }, dir)

        if len(value) == 0 {
          compilerErr("Expected some value as after ':'", dir, oper.Line)
        }

        hash[key] = value[0]
      }

      var ommVal OmmHash

      for k, v := range hash {
        ommVal.Set(k, v.Value)
      }

      return []Action{ Action{
        Type: "hash",
        Value: ommVal,
      } }
    case "[":

      var arr OmmArray

      for _, v := range item.Group {
        oper := makeOperations([][]Item{ v })[0]
        value := actionizer([]Operation{ oper }, dir)

        if len(value) == 0 {
          compilerErr("Each entry in the array must have a value", dir, oper.Line)
        }

        arr.PushBack(value[0].Value)
      }

      return []Action{ Action{
        Type: "array",
        Value: arr,
      } }
    case "expression value":

      var val = item.Token.Name

      if val[0] == '"' || val[0] == '`' { //detect string
        var str = OmmString{}
        str.FromGoType(val[1:len(val) - 1])
        return []Action{ Action{
          Type: "string",
          Value: str,
        } }
      } else if val[0] == '\'' { //detect a rune
        var oRune = OmmRune{}

        qrem := val[1:len(val) - 1] //remove quotes

        if len(qrem) != 1 {
          compilerErr("Runes can only be one character long", dir, item.Line)
        }

        oRune.FromGoType([]rune(qrem)[0])
        return []Action{ Action{
          Type: "rune",
          Value: oRune,
        } }
      } else if val == "true" || val == "false" { //detect a bool
        var boolean = OmmBool{}
        boolean.FromGoType(val == "true" /* convert to a boolean */)
        return []Action{ Action{
          Type: "bool",
          Value: boolean,
        } }
      } else if val == "undef" { //detect a falsey value
        var undef OmmUndef
        return []Action{ Action{
          Type: "falsey",
          Value: undef,
        } }
      } else if unicode.IsDigit(rune(val[0])) || val[0] == '.' || val[0] == '+' || val[0] == '-' { //detect a number
        var number = OmmNumber{}
        number.FromString(val)
        return []Action{ Action{
          Type: "number",
          Value: number,
        } }
      } else if val[0] == '$' { //detect a variable
        return []Action{ Action{
          Type: "variable",
          Name: val,
        } }
      } else { //detect nothing, which throws an error
        compilerErr(val + " is not a value", dir, item.Line)
      }

  }

  return []Action{}
}
