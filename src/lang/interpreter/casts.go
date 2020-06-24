package interpreter

import "strconv"

func cast(val Action, nType string) Action {
  valType := val.Type

  val.Type = nType

  switch valType + "->" + nType {
    case "string->number":
      val.Integer, val.Decimal = BigNumConverter(val.ExpStr)
    case "number->string":
      val.ExpStr = num_normalize(val)

      for k, v := range val.ExpStr {
        ommRune := emptyRune
        ommRune.Access = "public"
        ommRune.ExpStr = string(v)

        val.Hash_Values[strconv.Itoa(k)] = []Action{ ommRune }
      }
  }

  return val
}
